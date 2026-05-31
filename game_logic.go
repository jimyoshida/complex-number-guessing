package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	minValue   = 1
	maxValue   = 10
	maxGuesses = 7
)

type GuessResult int

const (
	TooLow GuessResult = iota
	TooHigh
	Correct
)

type AngleResult int

const (
	Left AngleResult = iota
	Right
	CorrectAngle
)

type ComparisonResult struct {
	Modulus GuessResult
	Angle   AngleResult
}

func parseComplexNumber(input string) (complex128, error) {
	input = strings.TrimSpace(input)

	c, err := strconv.ParseComplex("("+input+")", 128)
	if err == nil {
		return c, nil
	}

	r, err := strconv.ParseFloat(input, 64)
	if err == nil {
		return complex(r, 0), nil
	}

	return 0, fmt.Errorf("please enter a valid complex number (e.g., 5+3i, 5, or 3i)")
}

func validateGuess(guess complex128) error {
	re := real(guess)
	im := imag(guess)

	if re < minValue || re > maxValue || im < minValue || im > maxValue {
		return fmt.Errorf("please enter a complex number with both real and imaginary parts between %d and %d (e.g., 5+3i)", minValue, maxValue)
	}
	return nil
}

func getModulus(c complex128) float64 {
	return math.Sqrt(real(c)*real(c) + imag(c)*imag(c))
}

func getAngle(c complex128) float64 {
	return math.Atan2(imag(c), real(c))
}

func checkModulus(guess, target complex128) GuessResult {
	guessModulus := getModulus(guess)
	targetModulus := getModulus(target)

	const epsilon = 1e-9
	if math.Abs(guessModulus-targetModulus) < epsilon {
		return Correct
	}
	if guessModulus < targetModulus {
		return TooLow
	}
	return TooHigh
}

func checkAngle(guess, target complex128) AngleResult {
	guessAngle := getAngle(guess)
	targetAngle := getAngle(target)

	diff := targetAngle - guessAngle
	for diff > math.Pi {
		diff -= 2 * math.Pi
	}
	for diff <= -math.Pi {
		diff += 2 * math.Pi
	}

	const epsilon = 1e-9
	if math.Abs(diff) < epsilon {
		return CorrectAngle
	}
	if diff > 0 {
		return Left
	}
	return Right
}

func checkGuess(guess, target complex128) ComparisonResult {
	return ComparisonResult{
		Modulus: checkModulus(guess, target),
		Angle:   checkAngle(guess, target),
	}
}
