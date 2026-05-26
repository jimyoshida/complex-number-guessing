package main

import (
	"testing"
)

// TestParseComplexNumber tests parsing of complex number inputs in various formats.
func TestParseComplexNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid formats
		{"simple real", "5", false},
		{"positive complex", "5+3i", false},
		{"negative imaginary", "5-3i", false},
		{"pure imaginary", "3i", false},
		{"decimal real", "5.5+3i", false},
		{"decimal both", "5.5+3.5i", false},

		// Invalid formats
		{"invalid letters", "5x3i", true},
		{"missing i", "5+3", true},
		{"bad format", "five", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseComplexNumber(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseComplexNumber(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

// TestValidateGuess tests the input validation logic that ensures user guesses are within
// the valid range [minValue, maxValue] for both real and imaginary parts.
func TestValidateGuess(t *testing.T) {
	tests := []struct {
		name    string
		guess   complex128
		wantErr bool
	}{
		// Valid cases
		{"valid min", complex(float64(minValue), float64(minValue)), false},
		{"valid max", complex(float64(maxValue), float64(maxValue)), false},
		{"valid middle", complex(5, 5), false},

		// Invalid cases - real part out of range
		{"real below min", complex(0, 5), true},
		{"real above max", complex(11, 5), true},

		// Invalid cases - imaginary part out of range
		{"imag below min", complex(5, 0), true},
		{"imag above max", complex(5, 11), true},

		// Invalid cases - both out of range
		{"both below min", complex(0, 0), true},
		{"both above max", complex(11, 11), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGuess(tt.guess)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateGuess(%v) error = %v, wantErr %v", tt.guess, err, tt.wantErr)
			}
		})
	}
}

// TestCheckModulus tests the modulus comparison logic that determines if a guess
// has a magnitude (distance from origin) that's too low, too high, or correct.
func TestCheckModulus(t *testing.T) {
	tests := []struct {
		name     string
		guess    complex128
		target   complex128
		expected GuessResult
	}{
		// Winning condition - same modulus
		{"exact match", complex(5, 0), complex(5, 0), Correct},
		{"different but same modulus", complex(3, 4), complex(0, 5), Correct}, // both have modulus 5

		// Too low cases
		{"modulus too low", complex(2, 0), complex(4, 0), TooLow},
		{"modulus too low complex", complex(3, 4), complex(4, 4), TooLow}, // 5 < sqrt(32)

		// Too high cases
		{"modulus too high", complex(8, 0), complex(5, 0), TooHigh},
		{"modulus too high complex", complex(5, 5), complex(3, 4), TooHigh}, // sqrt(50) > 5
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkModulus(tt.guess, tt.target)
			if result != tt.expected {
				t.Errorf("checkModulus(%v, %v) = %v, want %v", tt.guess, tt.target, result, tt.expected)
			}
		})
	}
}

// TestCheckAngle tests the angle comparison logic that determines if a guess
// needs to turn left (counterclockwise), right (clockwise), or has the correct angle.
func TestCheckAngle(t *testing.T) {
	tests := []struct {
		name     string
		guess    complex128
		target   complex128
		expected AngleResult
	}{
		// Correct angle cases
		{"exact angle match", complex(50, 0), complex(50, 0), CorrectAngle},
		{"same angle different modulus", complex(3, 4), complex(6, 8), CorrectAngle}, // same angle 53.1°

		// Turn left (counterclockwise)
		{"turn left", complex(50, 0), complex(0, 50), Left},   // 0° to 90°
		{"turn left 2", complex(1, 0), complex(1, 1), Left},   // 0° to 45°

		// Turn right (clockwise)
		{"turn right", complex(0, 50), complex(50, 0), Right},  // 90° to 0°
		{"turn right 2", complex(1, 1), complex(1, 0), Right},  // 45° to 0°
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkAngle(tt.guess, tt.target)
			if result != tt.expected {
				t.Errorf("checkAngle(%v, %v) = %v, want %v", tt.guess, tt.target, result, tt.expected)
			}
		})
	}
}

// TestGameConstants verifies that the game configuration constants have the expected values.
// These constants define the game rules, so changing them should be intentional and tested.
// This test serves as documentation of the game's design:
// - minValue (1): Game uses complex numbers with parts starting from 1
// - maxValue (10): Game uses complex numbers with parts up to 10
// - maxGuesses (7): Player gets 7 attempts to guess
func TestGameConstants(t *testing.T) {
	if minValue != 1 {
		t.Errorf("minValue = %d, want 1", minValue)
	}
	if maxValue != 10 {
		t.Errorf("maxValue = %d, want 10", maxValue)
	}
	if maxGuesses != 7 {
		t.Errorf("maxGuesses = %d, want 7", maxGuesses)
	}
}

// TestGuessResultEnums verifies the enum values for GuessResult remain stable.
// The main game logic relies on these specific values, so this test ensures:
// - The enum is defined correctly
// - Future code changes don't accidentally alter the values
// - Consistency with any serialization or external interfaces
func TestGuessResultEnums(t *testing.T) {
	if TooLow != 0 {
		t.Errorf("TooLow = %v, want 0", TooLow)
	}
	if TooHigh != 1 {
		t.Errorf("TooHigh = %v, want 1", TooHigh)
	}
	if Correct != 2 {
		t.Errorf("Correct = %v, want 2", Correct)
	}
}
