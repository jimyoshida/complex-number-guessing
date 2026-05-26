package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
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

	// Try parsing as complex number directly
	c, err := strconv.ParseComplex("("+input+")", 128)
	if err == nil {
		return c, nil
	}

	// Try parsing as real number only
	real, err := strconv.ParseFloat(input, 64)
	if err == nil {
		return complex(real, 0), nil
	}

	return 0, fmt.Errorf("please enter a valid complex number (e.g., 5+3i, 5, or 3i)")
}

func validateGuess(guess complex128) error {
	real := real(guess)
	imag := imag(guess)

	if real < minValue || real > maxValue || imag < minValue || imag > maxValue {
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

	// Calculate the angular difference, normalized to (-π, π]
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

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	scanner := bufio.NewScanner(os.Stdin)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")). // Magenta
		Bold(true).
		Padding(1, 2)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")) // Cyan

	fmt.Println(titleStyle.Render("🎮 Complex Number Guessing Game"))
	fmt.Println(instructionStyle.Render(
		fmt.Sprintf("I'm thinking of a complex number between %d+%di and %d+%di.\n"+
			"Comparison is based on modulus. You have %d guesses. Good luck!\n"+
			"Enter guesses in format: 5+3i\n", minValue, minValue, maxValue, maxValue, maxGuesses)))
	fmt.Println()

	// Generate random complex number
	realPart := rng.Intn(maxValue-minValue+1) + minValue
	imagPart := rng.Intn(maxValue-minValue+1) + minValue
	target := complex(float64(realPart), float64(imagPart))
	guesses := 0

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")). // Yellow
		Bold(true)

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")). // Bright green
		Bold(true)
	modulusTooLowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")) // Red
	modulusTooHighStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("33")) // Blue
	modulusCorrectStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")) // Bright green
	angleTurnLeftStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")) // Pink/Magenta
	angleTurnRightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")) // Light blue
	angleCorrectStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")) // Bright green
	loseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")). // Bright red
		Bold(true)

	for guesses < maxGuesses {
		remaining := maxGuesses - guesses
		fmt.Print(promptStyle.Render(fmt.Sprintf("Guess %d/%d: ", guesses+1, maxGuesses)))

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		guess, err := parseComplexNumber(input)
		if err != nil {
			errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
			fmt.Println(errStyle.Render("❌ " + err.Error()))
			continue
		}

		if err := validateGuess(guess); err != nil {
			errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
			fmt.Println(errStyle.Render("❌ " + err.Error()))
			continue
		}

		guesses++

		result := checkGuess(guess, target)

		// Check if both modulus and angle are correct
		if result.Modulus == Correct && result.Angle == CorrectAngle {
			fmt.Println(successStyle.Render(fmt.Sprintf("\n✅ Correct! The number was %v.", target)))
			fmt.Println(successStyle.Render(fmt.Sprintf("🎉 You got it in %d guess(es)!", guesses)))
			return
		}

		// Show modulus feedback
		var modulusFeedback string
		var modulusResultStyle lipgloss.Style
		switch result.Modulus {
		case TooLow:
			modulusFeedback = fmt.Sprintf("📉 Modulus too low (%.2f)", getModulus(guess))
			modulusResultStyle = modulusTooLowStyle
		case TooHigh:
			modulusFeedback = fmt.Sprintf("📈 Modulus too high (%.2f)", getModulus(guess))
			modulusResultStyle = modulusTooHighStyle
		case Correct:
			modulusFeedback = "✓ Modulus correct"
			modulusResultStyle = modulusCorrectStyle
		}

		// Show angle feedback
		var angleFeedback string
		var angleResultStyle lipgloss.Style
		angleStr := fmt.Sprintf("(%.1f°)", math.Atan2(imag(guess), real(guess))*180/math.Pi)
		switch result.Angle {
		case Left:
			angleFeedback = fmt.Sprintf("↶ Turn left %s", angleStr)
			angleResultStyle = angleTurnLeftStyle
		case Right:
			angleFeedback = fmt.Sprintf("↷ Turn right %s", angleStr)
			angleResultStyle = angleTurnRightStyle
		case CorrectAngle:
			angleFeedback = fmt.Sprintf("✓ Angle correct %s", angleStr)
			angleResultStyle = angleCorrectStyle
		}

		fmt.Print(modulusResultStyle.Render(modulusFeedback) + " | " + angleResultStyle.Render(angleFeedback))

		if remaining-1 > 0 {
			fmt.Printf(" %d guess(es) remaining.\n", remaining-1)
		} else {
			fmt.Println()
		}
	}

	fmt.Println(loseStyle.Render(fmt.Sprintf("\n☠️  Out of guesses! The number was %v.", target)))

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
