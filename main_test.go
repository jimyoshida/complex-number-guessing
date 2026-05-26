package main

import (
	"testing"
)

// TestValidateGuess tests the input validation logic that ensures user guesses are within
// the valid range [minNumber, maxNumber]. This is critical because:
// - Invalid input can break game logic or provide a poor user experience
// - Boundary cases (min, max) are common sources of bugs
// - Out-of-range inputs should be rejected consistently
//
// Test coverage includes:
// - Valid boundary values (1, 100) to ensure min/max are accepted
// - Valid middle values (50) to ensure normal range is accepted
// - Invalid values just outside boundaries (0, 101) to catch off-by-one errors
// - Invalid values far outside range to ensure robust rejection
func TestValidateGuess(t *testing.T) {
	tests := []struct {
		name    string
		guess   int
		wantErr bool
	}{
		// Valid cases - should not return an error
		{"valid min", minNumber, false},
		{"valid max", maxNumber, false},
		{"valid middle", 50, false},

		// Invalid cases - should return an error
		{"below min", minNumber - 1, true},
		{"above max", maxNumber + 1, true},
		{"far below", 0, true},
		{"far above", 101, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGuess(tt.guess)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateGuess(%d) error = %v, wantErr %v", tt.guess, err, tt.wantErr)
			}
		})
	}
}

// TestCheckGuess tests the core game logic that compares user guesses against the target number.
// This function is the heart of the game and must correctly categorize every guess.
// Correctness here is essential for:
// - Giving accurate feedback to the player
// - Winning condition (when guess == target)
// - Game progression (tracking guesses)
//
// Test coverage includes:
// - Exact match (guess == target) to verify win condition
// - Too low guesses to ensure "lower" feedback is accurate
// - Too high guesses to ensure "higher" feedback is accurate
// - Boundary cases (1 vs 100) to catch edge-case logic errors
// - Off-by-one cases (49 vs 50, 51 vs 50) to verify comparison operators
// - Large differences to ensure logic works across full range
// - Negative numbers to verify robustness beyond game boundaries
func TestCheckGuess(t *testing.T) {
	tests := []struct {
		name     string
		guess    int
		target   int
		expected GuessResult
	}{
		// Winning condition - exact match
		{"guess equals target", 50, 50, Correct},

		// Standard too low cases
		{"guess below target", 25, 50, TooLow},
		{"guess 1 below target", 49, 50, TooLow},
		{"negative guess vs positive target", -5, 50, TooLow},
		{"large difference low", 10, 99, TooLow},
		{"guess min vs max", 1, 100, TooLow},

		// Standard too high cases
		{"guess above target", 75, 50, TooHigh},
		{"guess 1 above target", 51, 50, TooHigh},
		{"large difference high", 90, 10, TooHigh},
		{"guess max vs min", 100, 1, TooHigh},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkGuess(tt.guess, tt.target)
			if result != tt.expected {
				t.Errorf("checkGuess(%d, %d) = %v, want %v", tt.guess, tt.target, result, tt.expected)
			}
		})
	}
}

// TestGameConstants verifies that the game configuration constants have the expected values.
// These constants define the game rules, so changing them should be intentional and tested.
// This test serves as documentation of the game's design:
// - minNumber (1): Game asks for numbers starting from 1
// - maxNumber (100): Game asks for numbers up to 100
// - maxGuesses (7): Player gets 7 attempts to guess
func TestGameConstants(t *testing.T) {
	if minNumber != 1 {
		t.Errorf("minNumber = %d, want 1", minNumber)
	}
	if maxNumber != 100 {
		t.Errorf("maxNumber = %d, want 100", maxNumber)
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
