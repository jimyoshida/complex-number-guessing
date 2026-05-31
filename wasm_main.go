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

var (
	gameTarget  complex128
	gameGuesses int
	gameOver    bool
)

func jsNewGame(this js.Value, args []js.Value) any {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	gameTarget = complex(
		float64(rng.Intn(maxValue-minValue+1)+minValue),
		float64(rng.Intn(maxValue-minValue+1)+minValue),
	)
	gameGuesses = 0
	gameOver = false
	b, _ := json.Marshal(map[string]any{"maxGuesses": maxGuesses})
	return string(b)
}

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

	guess, err := parseComplexNumber(input)
	if err != nil {
		return errJSON(err.Error())
	}
	if err := validateGuess(guess); err != nil {
		return errJSON(err.Error())
	}

	gameGuesses++
	result := checkGuess(guess, gameTarget)

	won := result.Modulus == Correct && result.Angle == CorrectAngle
	gameOver = won || gameGuesses >= maxGuesses

	modulusStr := [3]string{"too_low", "too_high", "correct"}[result.Modulus]
	angleStr := [3]string{"left", "right", "correct"}[result.Angle]

	resp := map[string]any{
		"guessNumber":      gameGuesses,
		"guessesRemaining": maxGuesses - gameGuesses,
		"modulus": map[string]any{
			"result": modulusStr,
			"value":  math.Round(getModulus(guess)*100) / 100,
		},
		"angle": map[string]any{
			"result":  angleStr,
			"degrees": math.Round(getAngle(guess)*180/math.Pi*10) / 10,
		},
		"gameOver": gameOver,
		"won":      won,
		"target":   nil,
	}
	if gameOver {
		resp["target"] = fmt.Sprintf("%v", gameTarget)
	}

	b, _ := json.Marshal(resp)
	return string(b)
}

func main() {
	js.Global().Set("newGame", js.FuncOf(jsNewGame))
	js.Global().Set("makeGuess", js.FuncOf(jsMakeGuess))
	select {}
}
