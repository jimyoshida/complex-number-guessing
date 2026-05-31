//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
	"time"
)

// Game state shared across JS calls
var (
	gameTarget  complex128
	gameGuesses int
	gameOver    bool
)

// jsNewGame resets game state and generates a new random target.
// Returns JSON: {"maxGuesses": 7}
func jsNewGame(this js.Value, args []js.Value) any {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generate target with real and imaginary parts each in [minValue, maxValue]
	gameTarget = complex(
		float64(rng.Intn(maxValue-minValue+1)+minValue),
		float64(rng.Intn(maxValue-minValue+1)+minValue),
	)
	gameGuesses = 0
	gameOver = false
	b, _ := json.Marshal(map[string]any{"maxGuesses": maxGuesses})
	return string(b)
}

// jsMakeGuess validates and evaluates one guess against the current target.
// Invalid or out-of-range input does not consume a guess.
// Returns JSON with modulus/angle feedback, or {"error": "..."} on bad input.
func jsMakeGuess(this js.Value, args []js.Value) any {
	errJSON := func(msg string) string {
		b, _ := json.Marshal(map[string]any{"error": msg})
		return string(b)
	}

	if gameOver {
		return errJSON("game is over, start a new game")
	}
	if len(args) == 0 {
		return errJSON("no input provided")
	}

	input := args[0].String()

	// Parse and validate without consuming a guess
	guess, err := parseComplexNumber(input)
	if err != nil {
		return errJSON(err.Error())
	}
	if err := validateGuess(guess); err != nil {
		return errJSON(err.Error())
	}

	// Valid guess: increment counter and evaluate
	gameGuesses++
	result := checkGuess(guess, gameTarget)

	won := result.Modulus == Correct && result.Angle == CorrectAngle
	gameOver = won || gameGuesses >= maxGuesses

	// Map enum values to JSON-friendly strings
	modulusStr := [3]string{"too_low", "too_high", "correct"}[result.Modulus]
	angleStr := [3]string{"left", "right", "correct"}[result.Angle]

	resp := map[string]any{
		"guessNumber":      gameGuesses,
		"guessesRemaining": maxGuesses - gameGuesses,
		"modulus": map[string]any{
			"result": modulusStr,
			"value":  math.Round(getModulus(guess)*100) / 100,   // 2 decimal places
		},
		"angle": map[string]any{
			"result":  angleStr,
			"degrees": math.Round(getAngle(guess)*180/math.Pi*10) / 10, // 1 decimal place
		},
		"gameOver": gameOver,
		"won":      won,
		"target":   nil, // revealed only when the game ends
	}
	if gameOver {
		resp["target"] = fmt.Sprintf("%v", gameTarget)
	}

	b, _ := json.Marshal(resp)
	return string(b)
}

func main() {
	// Expose game functions to JavaScript on the global object
	js.Global().Set("newGame", js.FuncOf(jsNewGame))
	js.Global().Set("makeGuess", js.FuncOf(jsMakeGuess))
	select {} // block forever to keep the WASM instance alive
}
