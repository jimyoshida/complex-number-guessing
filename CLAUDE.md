# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A terminal-based complex number guessing game written in Go. The player has 7 attempts to guess a randomly generated complex number (real and imaginary parts each between 1 and 10). After each guess, feedback is provided on both the **modulus** (magnitude) and **angle** (argument) relative to the target.

The game also ships as a Web UI via Go WASM, served as a static site (e.g., GitHub Pages).

## Development Commands

All common tasks are in the `Makefile`:

| Command          | Description                                                   |
|------------------|---------------------------------------------------------------|
| `make build`     | Compile CLI binary to `bin/complex-number-guessing`          |
| `make run`       | Build then run the CLI binary                                 |
| `make test`      | Run all tests (`go test -v ./...`)                            |
| `make coverage`  | HTML coverage report at `build/coverage.html`                 |
| `make web-build` | Compile to WASM (`web/game.wasm`) and copy `wasm_exec.js`    |
| `make web-serve` | Build WASM then serve `web/` at `http://localhost:8080`       |
| `make clean`     | Remove `bin/`, `build/`, WASM artifacts, `coverage.out`       |

To run a single test: `go test -v -run TestCheckAngle ./...`

## Architecture

### Build tags separate two entry points

- [main.go](main.go) — `//go:build !wasm` — CLI entry point using `bufio`/`os` for I/O and `lipgloss` for terminal styling.
- [wasm_main.go](wasm_main.go) — `//go:build js && wasm` — WASM entry point; exposes `newGame` and `makeGuess` as JS globals via `syscall/js`. Holds game state in package-level vars (`gameTarget`, `gameGuesses`, `gameOver`).
- [game_logic.go](game_logic.go) — no build tag; shared by both builds. Contains all constants, types, and pure game logic functions.

### Shared game logic ([game_logic.go](game_logic.go))

- Constants: `minValue=1`, `maxValue=10`, `maxGuesses=7`
- Types: `GuessResult` (`TooLow/TooHigh/Correct`), `AngleResult` (`Left/Right/CorrectAngle`), `ComparisonResult`
- `parseComplexNumber` — wraps `strconv.ParseComplex` with a real-only fallback; accepts `5+3i`, `5`, `3i`, decimals
- `validateGuess` — checks both real and imaginary parts are within `[minValue, maxValue]`
- `checkModulus` / `checkAngle` / `checkGuess` — pure comparison functions; epsilon `1e-9` for float equality

### Web frontend ([web/](web/))

Vanilla JS + CSS. Calls `newGame()` and `makeGuess(input)` synchronously after WASM loads. `wasm_exec.js` (copied from `$(GOROOT)/misc/wasm/`) must be present alongside `game.wasm` — `make web-build` handles this.

## Key Invariants

- Invalid input (parse error or out-of-range) does **not** consume a guess — `continue` in the CLI loop, `{"error": "..."}` JSON in WASM, both without incrementing the counter.
- The target always has **both** real and imaginary parts (no pure-real targets).
- Angle feedback: `Left` means the guess angle is clockwise of the target (player should turn counterclockwise); `Right` is the opposite.
- `TooLow=0`, `TooHigh=1`, `Correct=2` — the WASM response maps these positionally via `[3]string{"too_low","too_high","correct"}[result.Modulus]`; don't reorder the iota.

## Testing Conventions

Table-driven tests with `t.Run` in [main_test.go](main_test.go). Tests cover `parseComplexNumber`, `validateGuess`, `checkModulus`, `checkAngle`, plus constants and enum value pins.

Run tests before committing: `make test`

## Commit Message Convention

Follow **Conventional Commits** (`feat:`, `fix:`, `docs:`, `test:`, `chore:`, `refactor:`, `build:`).
