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

type Card struct {
	Suit  string
	Value string
}

type Deck []Card

var suits = []string{"Hearts", "Diamonds", "Clubs", "Spades"}
var values = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

type Hand []Card

func newDeck() Deck {
	var deck Deck
	for _, suit := range suits {
		for _, value := range values {
			deck = append(deck, Card{Suit: suit, Value: value})
		}
	}
	return deck
}

func (deck Deck) shuffle() Deck {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

func (deck *Deck) drawCard() Card {
	card := (*deck)[0]
	*deck = (*deck)[1:]
	return card
}

func (h Hand) calculateValue() int {
	totalValue := 0
	aces := 0

	for _, card := range h {
		switch card.Value {
		case "A":
			aces++
			totalValue += 11
		case "K", "Q", "J":
			totalValue += 10
		default:
			val, _ := strconv.Atoi(card.Value)
			totalValue += val
		}
	}

	for aces > 0 && totalValue > 21 {
		totalValue -= 10
		aces--
	}

	return totalValue
}

func (h Hand) String() string {
	var cards []string
	for _, card := range h {
		cards = append(cards, fmt.Sprintf("%s of %s", card.Value, card.Suit))
	}
	return strings.Join(cards, ", ")
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	deck := newDeck().shuffle()
	playerHand := Hand{}
	dealerHand := Hand{}
	balance := 1000
	bet := 0

	fmt.Println("Welcome to Blackjack!")
	for {
		fmt.Printf("Current balance: $%d\n", balance)
		fmt.Print("Enter your bet (q to quit): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "q" {
			break
		}

		var err error
		bet, err = strconv.Atoi(input)
		if err != nil || bet <= 0 || bet > balance {
			fmt.Println("Invalid bet amount. Try again.")
			continue
		}

		balance -= bet

		playerHand = Hand{deck.drawCard(), deck.drawCard()}
		dealerHand = Hand{deck.drawCard(), deck.drawCard()}

		fmt.Println("Dealer's hand:", dealerHand[0], ", [HIDDEN]")
		fmt.Println("Your hand:", playerHand)

	playerTurn:
		for {
			fmt.Printf("Your hand (%d): %s\n", playerHand.calculateValue(), playerHand)
			fmt.Print("What will you do? (hit/stand/double/split): ")

			action, _ := reader.ReadString('\n')
			action = strings.TrimSpace(action)
			switch action {
			case "hit":
				playerHand = append(playerHand, deck.drawCard())
				if playerHand.calculateValue() > 21 {
					fmt.Printf("Busted! Your hand (%d): %s\n", playerHand.calculateValue(), playerHand)
					break playerTurn
				}
			case "stand":
				break playerTurn
			case "double":
				if len(playerHand) == 2 && balance >= bet {
					balance -= bet
					bet *= 2
					playerHand = append(playerHand, deck.drawCard())
					if playerHand.calculateValue() > 21 {
						fmt.Printf("Busted! Your hand (%d): %s\n", playerHand.calculateValue(), playerHand)
					}
					break playerTurn
				} else {
					fmt.Println("Cannot double: not enough balance or hand is not eligible.")
				}
			case "split":
				fmt.Println("Split not implemented.")
			default:
				fmt.Println("Invalid action. Please type hit, stand, double, or split.")
			}

			if len(deck) < 10 {
				deck = newDeck().shuffle()
			}
		}

		if playerHand.calculateValue() <= 21 {
			for dealerHand.calculateValue() < 17 {
				dealerHand = append(dealerHand, deck.drawCard())
			}
			fmt.Printf("Dealer's hand (%d): %s\n", dealerHand.calculateValue(), dealerHand)
		}

		playerValue := playerHand.calculateValue()
		dealerValue := dealerHand.calculateValue()
		switch {
		case playerValue > 21:
			fmt.Println("You busted and lost the bet.")
		case dealerValue > 21 || playerValue > dealerValue:
			fmt.Printf("You won $%d!\n", bet*2)
			balance += bet * 2
		case dealerValue == playerValue:
			fmt.Println("Push. You got your bet back.")
			balance += bet
		default:
			fmt.Println("Dealer wins. You lost the bet.")
		}

		fmt.Println()
	}

	fmt.Println("Thank you for playing! Your final balance is:", balance)
}
