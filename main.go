package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	minNumber = 1
	maxNumber = 100
	maxGuesses = 7
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Number Guessing Game ===")
	fmt.Printf("I'm thinking of a number between %d and %d.\n", minNumber, maxNumber)
	fmt.Printf("You have %d guesses. Good luck!\n\n", maxGuesses)

	target := rng.Intn(maxNumber-minNumber+1) + minNumber
	guesses := 0

	for guesses < maxGuesses {
		remaining := maxGuesses - guesses
		fmt.Printf("Guess %d/%d: ", guesses+1, maxGuesses)

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		guess, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Please enter a valid number.")
			continue
		}

		if guess < minNumber || guess > maxNumber {
			fmt.Printf("Please enter a number between %d and %d.\n", minNumber, maxNumber)
			continue
		}

		guesses++

		switch {
		case guess == target:
			fmt.Printf("\nCorrect! The number was %d.\n", target)
			fmt.Printf("You got it in %d guess(es)!\n", guesses)
			return
		case guess < target:
			fmt.Printf("Too low!")
		default:
			fmt.Printf("Too high!")
		}

		if remaining-1 > 0 {
			fmt.Printf(" %d guess(es) remaining.\n", remaining-1)
		} else {
			fmt.Println()
		}
	}

	fmt.Printf("\nOut of guesses! The number was %d.\n", target)
}
