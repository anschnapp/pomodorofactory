# PomodoroFactory

A terminal-based pomodoro timer that builds an ASCII art tomato in a virtual factory — complete with crane animations, welding sparks, and celebratory sounds.

## What It Does

PomodoroFactory runs a classic pomodoro cycle directly in your terminal:

- **25-minute work sessions** with a factory crane that welds together an ASCII tomato piece by piece
- **Automatic breaks** — 5 minutes after each pomodoro, 15 minutes after every 4th
- **Celebration sequence** when a pomodoro completes: colorful sparks, a fanfare, and Animalese-style gibberish speech reading a randomly generated congratulatory message
- **Progress tracking** with tomato emojis showing completed pomodoros in the current session

The entire UI is rendered with a custom zero-copy compositing engine — no TUI framework, just Go slices and ANSI escape codes.

## Requirements

- Go 1.24+
- A terminal that supports ANSI colors and alternate screen buffer
- Linux or macOS (for optional audio via `aplay` or `afplay`)

## Install & Run

```sh
go build -o pomodorofactory .
./pomodorofactory
```

Or run directly:

```sh
go run .
```

### Custom Duration

Pass a duration in minutes as the first argument (decimals allowed, useful for testing):

```sh
./pomodorofactory 0.2   # 12-second pomodoro
./pomodorofactory 50    # 50-minute pomodoro
```

## Controls

| Key | Action |
|-----|--------|
| `s` | Start a pomodoro (when idle) |
| `q` | Quit |
| `Ctrl+C` | Quit |

## How It Works

1. Press `s` to start — a factory crane begins welding the ASCII tomato row by row from the bottom up
2. The status bar shows a live countdown and your completed pomodoros
3. When the timer finishes, a two-phase celebration plays:
   - **Party**: colorful sparks overlay the tomato + rising tones and a fanfare
   - **Speech**: Animalese-style voice reads something like *"Spectacular we brilliantly forged a legendary pomodoro"*
4. A break starts automatically — short (5 min) or long (15 min) depending on your cycle
5. After the break, press `s` to start the next pomodoro

Audio is optional. If `aplay` (Linux) or `afplay` (macOS) isn't available, the celebration runs visual-only.

## Project Structure

```
main.go                  # Event loop, state machine, pomodoro cycle
pkg/
  factoryscene/          # Crane + welding animation
  motivationcloud/       # Rotating motivational phrases
  status/                # Status bar (countdown, tomato tracking)
  commandinput/          # Available keyboard actions display
  celebration/           # Two-phase completion ceremony
  audio/                 # PCM sound generation + playback (no deps)
  timer/                 # Countdown timer
  view/                  # Zero-copy slice compositor
  runecolor/             # Colored rune type for the rendering system
  render/                # Renderable interface
  slicehelper/           # Generic 2D slice utilities
  iohelper/              # String/file parsing helpers
```
