# CLAUDE.md

## Project Overview

A terminal-based complex number guessing game written in Go. The player is given 7 attempts to guess a randomly generated complex number (real and imaginary parts each between 1 and 10). After each guess, feedback is provided on both the **modulus** (magnitude) and the **angle** (argument) of the guess relative to the target.

## Repository Structure

```
.
├── main.go          # All game logic and the main entry point
├── main_test.go     # Unit tests for game logic functions
├── go.mod           # Go module definition (requires Go 1.21)
├── go.sum           # Dependency checksums
├── Makefile         # Build, run, test, and coverage targets
├── AGENTS.md        # Brief agent/AI conventions (superseded by this file)
└── .gitignore       # Ignores bin/, build/, vendor/, IDE files
```

Single-package project — everything lives in `package main`. There are no subdirectories or additional packages.

## Key Dependencies

- **`github.com/charmbracelet/lipgloss v0.9.1`** — terminal styling (colors, bold, padding). All visual output is rendered through lipgloss styles.
- Standard library: `bufio`, `fmt`, `math`, `math/rand`, `os`, `strconv`, `strings`, `time`.

## Development Commands

All common tasks are in the `Makefile`:

| Command          | Description                                        |
|------------------|----------------------------------------------------|
| `make build`     | Compile to `bin/complex-number-guessing`           |
| `make run`       | Build then run the binary                          |
| `make test`      | Run all tests with verbose output (`go test -v ./...`) |
| `make coverage`  | Generate HTML coverage report at `build/coverage.html` |
| `make clean`     | Remove `bin/`, `build/`, `coverage.out`            |

## Architecture & Game Logic

### Constants (`main.go:16-20`)
- `minValue = 1`, `maxValue = 10` — range for both real and imaginary parts of the target
- `maxGuesses = 7` — number of allowed attempts

### Core Types
- `GuessResult` (`TooLow`, `TooHigh`, `Correct`) — modulus comparison outcome
- `AngleResult` (`Left`, `Right`, `CorrectAngle`) — angle comparison outcome
- `ComparisonResult` — wraps both results for a single guess

### Key Functions
| Function | Location | Purpose |
|---|---|---|
| `parseComplexNumber` | `main.go:43` | Parses user input into `complex128`; accepts `5+3i`, `5`, `3i`, decimals |
| `validateGuess` | `main.go:61` | Ensures both parts are within `[minValue, maxValue]` |
| `checkModulus` | `main.go:79` | Compares modulus (|guess| vs |target|) with epsilon `1e-9` |
| `checkAngle` | `main.go:93` | Compares angle using `math.Atan2`, normalizes difference to `(-π, π]` |
| `checkGuess` | `main.go:116` | Combines modulus and angle checks into `ComparisonResult` |

### Input Handling
Invalid input (parse error or out-of-range) does **not** consume a guess — the loop uses `continue` and the guess counter only increments after a valid, in-range input.

## Testing Conventions

Tests are in `main_test.go` as table-driven tests using `t.Run`. Each test function targets one exported-or-internal function:

- `TestParseComplexNumber` — valid and invalid input format coverage
- `TestValidateGuess` — boundary value testing (min, max, out-of-range)
- `TestCheckModulus` — same-modulus (different numbers), too low, too high
- `TestCheckAngle` — same angle (different modulus), left/right turns
- `TestGameConstants` — documents and pins `minValue`, `maxValue`, `maxGuesses`
- `TestGuessResultEnums` — pins `TooLow=0`, `TooHigh=1`, `Correct=2`

Run tests before committing: `make test`

## Commit Message Convention

Follow **Conventional Commits** (`feat:`, `fix:`, `docs:`, `test:`, `chore:`, `refactor:`, `build:`). Examples from history:

```
feat: implement complex number guessing game with modulus and angle comparison
test: add comprehensive unit tests with testable game logic
build: add coverage target to Makefile
docs: add color code comments to lipgloss styles
```

## Styling Conventions

All terminal output uses `lipgloss.NewStyle()`. Color assignments:

| Style variable          | Color code | Usage                        |
|-------------------------|------------|------------------------------|
| `titleStyle`            | 205 (magenta) | Game title                |
| `instructionStyle`      | 86 (cyan)  | Opening instructions          |
| `promptStyle`           | 226 (yellow) | Guess prompt               |
| `successStyle`          | 46 (green) | Win message                   |
| `modulusTooLowStyle`    | 196 (red)  | Modulus feedback: too low     |
| `modulusTooHighStyle`   | 33 (blue)  | Modulus feedback: too high    |
| `modulusCorrectStyle`   | 46 (green) | Modulus feedback: correct     |
| `angleTurnLeftStyle`    | 205 (pink) | Angle feedback: turn left     |
| `angleTurnRightStyle`   | 117 (light blue) | Angle feedback: turn right |
| `angleCorrectStyle`     | 46 (green) | Angle feedback: correct       |
| `loseStyle`             | 196 (red)  | Lose message                  |

## Important Notes

- The game generates a target with **both** real and imaginary parts, so pure-real targets never occur.
- Feedback for angle direction: `Left` means the guess angle is **clockwise** of the target (turn counterclockwise to get there); `Right` means counterclockwise of the target.
- Epsilon comparisons (`1e-9`) are used for both modulus and angle checks to handle floating-point precision.
- Build artifacts (`bin/`, `build/`) and `coverage.out` are gitignored — never commit them.
