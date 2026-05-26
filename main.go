package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	minNumber = 1
	maxNumber = 100
	maxGuesses = 7
)

type GuessResult int

const (
	TooLow GuessResult = iota
	TooHigh
	Correct
)

func validateGuess(guess int) error {
	if guess < minNumber || guess > maxNumber {
		return fmt.Errorf("please enter a number between %d and %d", minNumber, maxNumber)
	}
	return nil
}

func checkGuess(guess, target int) GuessResult {
	if guess == target {
		return Correct
	}
	if guess < target {
		return TooLow
	}
	return TooHigh
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

	fmt.Println(titleStyle.Render("🎮 Number Guessing Game"))
	fmt.Println(instructionStyle.Render(
		fmt.Sprintf("I'm thinking of a number between %d and %d.\n"+
			"You have %d guesses. Good luck!\n", minNumber, maxNumber, maxGuesses)))
	fmt.Println()

	target := rng.Intn(maxNumber-minNumber+1) + minNumber
	guesses := 0

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")). // Yellow
		Bold(true)

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")). // Bright green
		Bold(true)
	lowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("208")) // Orange
	highStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("33")) // Blue
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
		guess, err := strconv.Atoi(input)
		if err != nil {
			errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
			fmt.Println(errStyle.Render("❌ Please enter a valid number."))
			continue
		}

		if err := validateGuess(guess); err != nil {
			errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
			fmt.Println(errStyle.Render("❌ " + err.Error()))
			continue
		}

		guesses++

		result := checkGuess(guess, target)
		switch result {
		case Correct:
			fmt.Println(successStyle.Render(fmt.Sprintf("\n✅ Correct! The number was %d.", target)))
			fmt.Println(successStyle.Render(fmt.Sprintf("🎉 You got it in %d guess(es)!", guesses)))
			return
		case TooLow:
			fmt.Print(lowStyle.Render("📉 Too low!"))
		case TooHigh:
			fmt.Print(highStyle.Render("📈 Too high!"))
		}

		if remaining-1 > 0 {
			fmt.Printf(" %d guess(es) remaining.\n", remaining-1)
		} else {
			fmt.Println()
		}
	}

	fmt.Println(loseStyle.Render(fmt.Sprintf("\n☠️  Out of guesses! The number was %d.", target)))

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
