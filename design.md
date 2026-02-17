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
│  ·  │ (factory     │ (motivation cloud)   │  ·   │
│  ·  │  scene:      │                      │  ·   │
│  ·  │  crane+art)  │                      │  ·   │
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
| Factory Scene | `factoryscene` | Crane + welding animation building ASCII art | Dynamic: crane pillar on left, arm extends to weld point, flickering yellow sparks, art reveals L→R per row (bottom-to-top). Canvas dimensions fixed at construction to the max across all products. `MakeFactoryScene(products)` sizes the canvas; `LoadArt(art)` hot-swaps the target art (recomputes body rows/cells, resets progress) without resizing. `SetProgress(float64)` driven by `timer.Progress()`. Width = pillarWidth(1) + craneOverhead(4) + maxArtWidth. |
| Motivation Cloud | `motivationcloud` | Inspirational phrases | Dynamic: 5 phrases from a pool of 152 (8 categories). Every 15s one phrase is replaced with an animated transition — old phrase fades out char-by-char (right→left), new phrase reveals in char-by-char (left→right) with a dim leading edge. ~1.5s transition per swap. 3-color palette (HiCyan, White, HiMagenta). Animates in all states including idle. |
| Status | `status` | Pomodoro state info | Dynamic: shows state text + countdown (MM:SS) on line 1, achievement emojis for completed products on line 2. `SetAchievements(line1, emojis []string)` updates both; emojis joined with spaces. `statusWidth` = 50. |
| Command Input | `commandinput` | Available keyboard actions | Dynamic: fixed height=4, width=50. `SetTexts(commandText, selectorText)` fills two padded content lines between separator rows. In idle shows command line + selector (`build next:  ← [Name] →`); in other states selector line is blank. |
| ~~Pomodoro~~ | `pomodorobuild` | ~~ASCII art tomato with fill animation~~ | **Replaced** by `factoryscene`. Still in repo but unused by main. |

All four visible components update dynamically during the session.

| Audio Engine | `audio` | Programmatic sound generation + playback | Generates PCM samples (sine waves, noise, sawtooth) with pure Go math. Plays via `aplay` (Linux) or `afplay` (macOS, temp WAV file). No Go audio dependencies. |
| Celebration | `celebration` | Two-phase completion ceremony | State machine: PhaseNone → PhaseParty → PhaseSpeech → PhaseDone. `Start(message)` accepts a custom congratulatory message for the speech phase. Coordinates audio playback with TUI animation. |

## Rendering Pipeline (Current)

```
main()
  │
  ├─ put terminal in raw mode (golang.org/x/term)
  ├─ enter alternate screen buffer (ESC[?1049h)
  ├─ hide cursor (ESC[?25l)
  ├─ construct 4 Renderables
  │    └─ factoryscene.MakeFactoryScene(product.All)  // canvas sized to max product dims
  ├─ MakeView(topLeft, topRight, middle, bottom)
  │    ├─ calculate total layout dimensions (with 2v/5h margins)
  │    ├─ allocate master canvas with border
  │    └─ extract slice sub-regions for each component
  ├─ view.Render() + view.Print()   // initial frame
  └─ event loop:
       ├─ goroutine reads stdin as rune channel; ESC sequences decoded:
       │    ESC [ C → keyRight (rune(-2)), ESC [ D → keyLeft (rune(-1))
       ├─ 50ms ticker drives animation updates
       ├─ 'q' or Ctrl+C (0x03) → exit (defers restore terminal + leave alt screen)
       ├─ 'h' / ← or 'l' / → → cycle selected product (idle only)
       ├─ 's' → LoadArt(selectedProduct.Art) + start timer (idle only)
       ├─ on tick: update progress (float64) + status text → re-render + re-print
       ├─ on timer finish: play notification bell → enter waitingForCelebration state
       ├─ 'c' key: trigger full celebration (party sparks + sounds → gibberish speech)
       └─ state machine: idle → working → waitingForCelebration → celebrating → onBreak → idle (repeats)
```

Raw mode is needed so keypresses arrive immediately without Enter. Print uses `\r\n` because raw mode disables the kernel's `\n` → `\r\n` translation.

## Layout System

The View uses a fixed 4-slot layout with margins (2 vertical, 5 horizontal):

- **Top row**: two components side by side (pomodoro + motivation cloud)
- **Middle row**: one full-width component (status)
- **Bottom row**: one full-width component (command input)

Width is `max(top_combined, middle, bottom)`. Height is the sum of all rows plus margins.

## Products & ASCII Art

All buildable products live in `pkg/product/`. Each `Product` has a `Name`, `Emoji`, and `Art [][]runecolor.ColoredRune` (pre-colored, ready for factoryscene). `product.All` is the ordered registry, initialized via `init()`.

Three products ship:
| Product | Emoji | Color scheme |
|---|---|---|
| Tomato | 🍅 | `\|/\` → FgGreen, fill → FgRed |
| Coffee Cup | ☕ | `\|_-=` → FgHiYellow, `~` → FgHiWhite, fill → FgYellow |
| Penguin | 🐧 | `\|/\_^` → FgHiCyan, `o` → FgHiWhite, fill → FgHiBlack |

Art files are embedded via `//go:embed art/*.txt`, parsed by `iohelper.SplitMultilineStringToSlice`, and colored with per-character rune maps. The canvas in `factoryscene` is sized to the widest/tallest art at startup (currently the tomato at 23 cols × 10 rows), so all products fit without resizing the view region.

## Color System

Two approaches coexist:
- **fatih/color attributes**: standard named colors (FgGreen, FgRed) for component content
- **Raw SGR sequences**: used for the border background (RGB 100,100,100 via attribute codes `48, 2, R, G, B`)

Both are stored in `ColoredRune.ColorAttributes` and applied identically via `color.Set()` during printing.

## What's Missing (for a functional pomodoro app)

### 1. ~~Event Loop & Input Handling~~ ✓ Done
Terminal is in raw mode via `golang.org/x/term`. Goroutine reads stdin, event loop dispatches keypresses. `q` and Ctrl+C quit cleanly. Alternate screen buffer keeps the host terminal clean.

### 2. ~~Timer with Fill Animation~~ ✓ Done
`pkg/timer` provides a countdown timer (configurable via CLI arg in minutes, e.g. `0.1` for 6s). `timer.Progress()` returns fine-grained float64 (0.0–1.0). `status.SetText()` shows live MM:SS countdown. Event loop uses 50ms ticker + `'s'` key to start.

### 2b. ~~Factory Crane + Welding Animation~~ ✓ Done
`pkg/factoryscene` replaces `pomodorobuild` in the top-left slot. Combines a vertical crane pillar, horizontal arm, flickering welding sparks (bright yellow), and the ASCII art being built. Art reveals left-to-right per row, bottom-to-top row order. Each row gets equal time regardless of width (narrow rows = slower per-char, wide rows = faster per-char). The crane arm extends from the pillar through leading whitespace to the weld point; sparks sit at the left edge of content with a 1-space gap before the first revealed char. `contentOffset = pillarWidth(1) + craneOverhead(4)` guarantees room for arm/sparks/gap even on widest rows (firstCol=0).

### 2c. ~~Celebration on Completion~~ ✓ Done
Two-phase celebration triggers when the pomodoro timer finishes:

**Phase 1 — Party**: Factory scene overlays colorful sparks (yellow, green, magenta, cyan, red) on ~15% of the completed tomato art, randomly changing each tick. Status text flashes "POMODORO COMPLETE!" cycling through bright colors. Party sounds play: 3 rising sine sweeps + 2 noise-burst pops + a square-wave C-E-G-C fanfare.

**Phase 2 — Gibberish Speech**: Animalese-style voice reads a randomly generated congratulatory message. The message is composed from 4 word lists (20 words each): `[congrats] we [adverb] [verb] a [adjective] <product name>` — yielding 160,000 possible combinations per product. Each character maps to a short pitched blip (vowels: 200-400Hz/80ms, consonants: 400-800Hz/60ms, spaces: silence). Waveform is 70% sawtooth + 30% sine with ±15% random pitch variation per character. Status text shows the message with character-by-character highlight (spoken=white, current=bold yellow, upcoming=dim).

**Audio engine** (`pkg/audio/`): All sounds generated with pure Go math — no audio files or Go audio libraries. Platform-native playback: `aplay` on Linux (raw PCM via stdin), `afplay` on macOS (temp WAV file with 44-byte header). Audio is optional — if no playback tool is found, celebration runs visual-only. `statusWidth` bumped from 30→50 to fit the longer randomized messages.

### 3. ~~State Machine~~ ✓ Done
Full pomodoro cycle implemented in `main.go` with 5 states: `stateIdle` → `stateWorking` → `stateWaitingForCelebration` → `stateCelebrating` → `stateOnBreak` → back to `stateIdle`.

- **Cycle**: 4 pomodoros per set. Short break (5min) after pomodoros 1–3, long break (15min) after the 4th. Cycle repeats indefinitely.
- **User-triggered celebration**: When the timer finishes, the app enters `stateWaitingForCelebration` — a mechanical bell notification plays (`pkg/audio/notification.go`, 1.2s brrrrr), status shows "Pomodoro done! Press [c] to celebrate", and the command bar shows `[c]elebrate`. The user presses `c` whenever they're ready, which triggers the full celebration. This avoids interrupting the user mid-task with a forced animation.
- **Auto-break**: Break starts automatically after the celebration finishes — no keypress needed.
- **Factory wording**: Status uses factory-themed language ("Factory running", "Factory needs a short cooldown", "Factory needs a longer cooldown", "Factory ready").
- **Achievement tracking**: `achievedEmojis []string` collects each completed product's emoji in order. Displayed on status line 2 via `SetAchievements()`. Break duration is determined by `len(achievedEmojis) % pomodorosPerSet`.
- **Factory reset**: `factoryscene.Reset()` clears progress to 0 when break ends, so the next build starts fresh. During break, the completed art stays visible.
- **Timer reuse**: `timer.Reset(duration)` allows switching between work and break durations without creating a new timer.
- **Command input**: Dynamic via `commandinput.SetTexts(commandText, selectorText)` — idle shows `[s]tart | [q]uit` + selector row; working/break shows `[q]uit` with blank selector; waiting shows `[c]elebrate`.

### 4. ~~Dynamic Motivation Cloud~~ ✓ Done
152 phrases across 8 thematic categories (Focus, Encouragement, Progress, Energy, Mindset, Calm & Steady, Fun & Playful). 5 phrases displayed at a time, scattered across 10 rows with random indentation and color (3-color palette: HiCyan, White, HiMagenta).

Every 15 seconds, one random phrase is replaced with an animated transition: the old phrase fades out character-by-character from right to left, then the new phrase reveals in left to right — each with a dim leading/trailing edge. At 50ms per character, a typical 15-char phrase transitions in ~1.5s total. `ReplaceOne()` initiates the swap, `Tick()` advances animation each frame. Animates continuously in all states (idle, working, break).

### 5. ~~Multi-Product Factory~~ ✓ Done
The factory can build three products selectable before each pomodoro: Tomato 🍅 (default), Coffee Cup ☕, Penguin 🐧. Each product lives in `pkg/product/` with its own embedded ASCII art file and color scheme.

- **Selection**: In idle state, `h`/`←` and `l`/`→` cycle through `product.All`. The command bar's selector row shows `build next:  ← [Name] →`. Selection persists across breaks until the app exits.
- **Canvas stability**: `factoryscene.MakeFactoryScene(product.All)` sizes the canvas once to the largest product dimensions. `LoadArt(art)` hot-swaps the art at start time, recomputing body rows and cell positions without touching the view region.
- **Arrow key decoding**: `inputCh` is now `chan rune`. The stdin goroutine detects ESC sequences (`ESC [ C/D`) and sends sentinel values `keyRight = rune(-2)` / `keyLeft = rune(-1)`.
- **Achievement emojis**: Each completed product appends its `Emoji` to `achievedEmojis []string`, shown on status line 2 (e.g. `☕ 🐧 🍅`). The congratulatory speech ends with the product name instead of a fixed "pomodoro".

### 6. Deliberately Out of Scope
Task tracking and persistence were considered and intentionally skipped. The app is a focused pomodoro timer — task management belongs in the user's own system. Adding a task list would require significant UI rework and push the app toward being a todo manager.

## Utility Code Notes

- `iohelper.ReadFileInArray()` has a bug: it indexes into an empty slice (`lines[i]` where `lines` was initialized as `[][]rune{}`). This function is unused - embedding replaced file reading.
- `iohelper.SplitMultilineStringToSlice()` is the active helper, used for parsing embedded ASCII art.
- `slicehelper` uses Go generics (`[T any]`) for reusable 2D slice operations.
- `view.max()` is a custom variadic max function (predates Go 1.21's `max` builtin).

## Dependencies

- `github.com/fatih/color` (v1.18.0) - ANSI color output
- `golang.org/x/term` (v0.40.0) - terminal raw mode

The project deliberately avoids TUI frameworks, building its own rendering from scratch.
