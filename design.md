# PomodoroFactory - Architecture & Design

## Overview

PomodoroFactory is a terminal-based pomodoro timer written in Go with a custom rendering engine. It renders colored ASCII art components into a shared 2D buffer using Go's slice mechanics for zero-copy compositing.

## Core Architectural Idea: Shared-Array Slice Rendering

The central design insight is exploiting how Go slices work: a slice created from another slice shares the same underlying array. This is used to build a simple but effective compositing system:

1. A single master canvas (`completeView`) is allocated as `[][]runecolor.ColoredRune`
2. For each UI component, a **sub-region** is extracted using slice expressions:
   ```go
   viewRegion[row] = completeView[absoluteRow][startCol : startCol+width]
   ```
3. Each component receives its sub-region and writes into it via `Render(viewRegion)`
4. Because the sub-region IS the master canvas (same backing array), no copying or merging step is needed
5. Printing just iterates the master canvas once

This means: **components render independently, but their output lands directly in the final frame buffer.**

```
completeView (master canvas)
┌──────────────────────────────────────────────────┐
│  ·····border (RGB 100,100,100 background)······  │
│  ·  ┌─────────────┬──────────────────────┐  ·   │
│  ·  │ slice A      │ slice B              │  ·   │
│  ·  │ (pomodoro)   │ (motivation cloud)   │  ·   │
│  ·  └─────────────┴──────────────────────┘  ·   │
│  ·  ┌────────────────────────────────────┐  ·   │
│  ·  │ slice C (status)                   │  ·   │
│  ·  └────────────────────────────────────┘  ·   │
│  ·  ┌────────────────────────────────────┐  ·   │
│  ·  │ slice D (command input)            │  ·   │
│  ·  └────────────────────────────────────┘  ·   │
│  ··············································  │
└──────────────────────────────────────────────────┘

Slices A-D are NOT copies. They are views into completeView's memory.
Writing to slice A[2][5] actually writes to completeView[row][col].
```

## Key Types

### `runecolor.ColoredRune`
The atomic unit of the rendering system. Every cell in the canvas is one of these:
```go
type ColoredRune struct {
    Symbol          rune
    ColorAttributes []color.Attribute  // fatih/color attributes (incl. raw SGR)
}
```
Conversion helpers exist to map plain `[]rune` into `[]ColoredRune` with per-character color rules (e.g. `|` -> green, `x` -> red).

### `render.Renderable` Interface
Every UI component implements this:
```go
type Renderable interface {
    Render([][]runecolor.ColoredRune)  // write content into the given view region
    Width() int
    Height() int
}
```
Width/Height are used by the View to allocate the right sub-region. Render receives a slice into the master canvas.

### `view.View`
The compositor. Owns the master canvas and a list of `viewRegionRenderableBundle` entries (each pairing a Renderable with its slice region). Orchestrates layout, rendering, and printing.

### `slicehelper`
Generic utilities for 2D slices: `Copy2DSlice[T]`, `MaxWidth[T]`, `MinWidth[T]`. Used by components to copy their internal content into their assigned view region.

## Current Component Inventory

| Component | Package | Role | Status |
|---|---|---|---|
| Pomodoro | `pomodorobuild` | ASCII art tomato with color | Renders static art. Has `percentage` field (unused). |
| Motivation Cloud | `motivationcloud` | Inspirational phrases | Static placeholder. Plans for word lists + random selection. |
| Status | `status` | Pomodoro state info | Static placeholder ("Pomodoro running", "Finished: 3"). |
| Command Input | `commandinput` | Available keyboard actions | Static placeholder ("[s]tart \| [q]uit"). |

All components currently pre-build their content at construction time and copy it into the view region on `Render()`. There is no dynamic re-rendering yet.

## Rendering Pipeline (Current)

```
main()
  │
  ├─ put terminal in raw mode (golang.org/x/term)
  ├─ enter alternate screen buffer (ESC[?1049h)
  ├─ hide cursor (ESC[?25l)
  ├─ construct 4 Renderables (static content built in constructors)
  ├─ MakeView(topLeft, topRight, middle, bottom)
  │    ├─ calculate total layout dimensions (with 5px margins)
  │    ├─ allocate master canvas with border
  │    └─ extract slice sub-regions for each component
  ├─ view.Render() + view.Print()   // initial frame
  └─ event loop:
       ├─ goroutine reads stdin byte-by-byte → sends on channel
       ├─ blocks on input channel
       ├─ 'q' or Ctrl+C (0x03) → exit (defers restore terminal + leave alt screen)
       └─ other keys → re-render + re-print (cursor-home ESC[H before each frame)
```

Raw mode is needed so keypresses arrive immediately without Enter. Print uses `\r\n` because raw mode disables the kernel's `\n` → `\r\n` translation.

## Layout System

The View uses a fixed 4-slot layout with uniform margins (5 cells on all sides):

- **Top row**: two components side by side (pomodoro + motivation cloud)
- **Middle row**: one full-width component (status)
- **Bottom row**: one full-width component (command input)

Width is `max(top_combined, middle, bottom)`. Height is the sum of all rows plus margins.

## ASCII Art Embedding

The pomodoro art is loaded via Go's `//go:embed` directive from the `pomodoro-asci` file, then parsed into `[][]rune` by `iohelper.SplitMultilineStringToSlice`. Color is applied per-character using a rune-to-color map (structural chars like `|/\` get green, fill chars like `x` get red).

## Color System

Two approaches coexist:
- **fatih/color attributes**: standard named colors (FgGreen, FgRed) for component content
- **Raw SGR sequences**: used for the border background (RGB 100,100,100 via attribute codes `48, 2, R, G, B`)

Both are stored in `ColoredRune.ColorAttributes` and applied identically via `color.Set()` during printing.

## What's Missing (for a functional pomodoro app)

### 1. ~~Event Loop & Input Handling~~ ✓ Done
Terminal is in raw mode via `golang.org/x/term`. Goroutine reads stdin, event loop dispatches keypresses. `q` and Ctrl+C quit cleanly. Alternate screen buffer keeps the host terminal clean.

### 2. Timer Logic
The pomodoro timer itself doesn't exist yet. Needed:
- Countdown timer (configurable work/break durations)
- State machine: idle -> working -> short break -> working -> ... -> long break
- The `percentage` field on pomodorobuild hints at planned fill-level animation

### 3. Dynamic Component Updates
All components build their content once at construction. For a live app:
- Components need to update their internal state and re-render
- The Renderable interface may need an `Update(state)` method or components need references to shared state
- The motivation cloud should pick random phrases from word lists
- Status should reflect actual timer state and pomodoro count
- Command input should show context-appropriate actions

### 4. Pomodoro Visual Animation
The `percentage` field suggests the tomato should fill up as time passes (e.g. `x` chars appearing row by row, or color changing from empty to filled). This would be the signature visual feature.

### 5. State Persistence
From the ui-draft and intended features:
- Track completed pomodoros per day
- Task list with pomodoro counts per task
- This implies some form of storage (file-based likely, given the terminal nature)
- State format and storage location TBD - this can grow complex

### 6. Motivation Cloud
Intended as a rotating display of motivational phrases/words. Needs:
- Word/phrase lists (possibly embedded or configurable)
- Random or rotating selection
- Potentially state-aware (different phrases for work vs. break)

### 7. UI-Draft Features Not Yet Represented
From the `ui-draft` file, additional planned elements:
- Progress bar (vertical `|||` bars showing elapsed time)
- Today's pomodoro count with tomato emoji
- Task list with per-task pomodoro tracking (e.g. "merge roles and permissions 🍅🍅🍅🍅")

## Utility Code Notes

- `iohelper.ReadFileInArray()` has a bug: it indexes into an empty slice (`lines[i]` where `lines` was initialized as `[][]rune{}`). This function is unused - embedding replaced file reading.
- `iohelper.SplitMultilineStringToSlice()` is the active helper, used for parsing embedded ASCII art.
- `slicehelper` uses Go generics (`[T any]`) for reusable 2D slice operations.
- `view.max()` is a custom variadic max function (predates Go 1.21's `max` builtin).

## Dependencies

- `github.com/fatih/color` (v1.18.0) - ANSI color output
- `golang.org/x/term` (v0.40.0) - terminal raw mode

The project deliberately avoids TUI frameworks, building its own rendering from scratch.
