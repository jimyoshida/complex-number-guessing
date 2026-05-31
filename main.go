//go:build !wasm

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	scanner := bufio.NewScanner(os.Stdin)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(1, 2)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	fmt.Println(titleStyle.Render("🎮 Complex Number Guessing Game"))
	fmt.Println(instructionStyle.Render(
		fmt.Sprintf("I'm thinking of a complex number between %d+%di and %d+%di.\n"+
			"Comparison is based on modulus. You have %d guesses. Good luck!\n"+
			"Enter guesses in format: 5+3i\n", minValue, minValue, maxValue, maxValue, maxGuesses)))
	fmt.Println()

	realPart := rng.Intn(maxValue-minValue+1) + minValue
	imagPart := rng.Intn(maxValue-minValue+1) + minValue
	target := complex(float64(realPart), float64(imagPart))
	guesses := 0

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true)

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true)
	modulusTooLowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))
	modulusTooHighStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("33"))
	modulusCorrectStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46"))
	angleTurnLeftStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))
	angleTurnRightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117"))
	angleCorrectStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46"))
	loseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
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

		if result.Modulus == Correct && result.Angle == CorrectAngle {
			fmt.Println(successStyle.Render(fmt.Sprintf("\n✅ Correct! The number was %v.", target)))
			fmt.Println(successStyle.Render(fmt.Sprintf("🎉 You got it in %d guess(es)!", guesses)))
			return
		}

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
