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

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator using current time
	targetNumber := rand.Intn(100) + 1 // Random number between 1 and 100
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Guess the number between 1 and 100.")

	for {
		fmt.Print("Enter your guess: ")
		scanner.Scan() // Read input from user
		input := scanner.Text()
		guess, err := strconv.Atoi(strings.TrimSpace(input))

		if err != nil {
			fmt.Println("Please enter a valid number!")
			continue
		}

		if guess < targetNumber {
			fmt.Println("Too low! Try again.")
		} else if guess > targetNumber {
			fmt.Println("Too high! Try again.")
		} else {
			fmt.Printf("Congratulations! You guessed the right number: %d\n", targetNumber)
			break
		}
	}
}
