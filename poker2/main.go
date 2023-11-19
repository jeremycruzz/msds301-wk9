package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	Hearts   = "\u2665"
	Diamonds = "\u2666"
	Clubs    = "\u2663"
	Spades   = "\u2660"
)

type Suit string
type Rank string

const (
	Two   Rank = "2"
	Three Rank = "3"
	Four  Rank = "4"
	Five  Rank = "5"
	Six   Rank = "6"
	Seven Rank = "7"
	Eight Rank = "8"
	Nine  Rank = "9"
	Ten   Rank = "T"
	Jack  Rank = "J"
	Queen Rank = "Q"
	King  Rank = "K"
	Ace   Rank = "A"
)

var ranks = []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
var suits = []Suit{Hearts, Diamonds, Clubs, Spades}

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	return string(c.Rank) + string(c.Suit)
}

type Deck []Card

func (d Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
}

func NewDeck() Deck {
	var deck Deck
	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{Suit: suit, Rank: rank})
		}
	}
	return deck
}

func (d *Deck) Draw(num int) []Card {
	cards := (*d)[:num]
	*d = (*d)[num:]
	return cards
}

type Hand []Card

func (h Hand) String() string {
	sort.Slice(h, func(i, j int) bool {
		return h[i].Rank < h[j].Rank
	})
	var cards []string
	for _, card := range h {
		cards = append(cards, card.String())
	}
	return strings.Join(cards, " ")
}

func (h Hand) Score() (score int, highCard Card) {
	sort.Slice(h, func(i, j int) bool {
		return h[i].Rank < h[j].Rank
	})

	// Count occurrences of each rank
	rankCount := make(map[Rank]int)
	for _, card := range h {
		rankCount[card.Rank]++
	}

	// Straight, Flush, and Straight Flush check
	isFlush := true
	isStraight := true
	previousRankIndex := -1
	for i, card := range h {
		if i > 0 {
			if card.Suit != h[i-1].Suit {
				isFlush = false
			}
			currentRankIndex := indexOfRank(card.Rank)
			if currentRankIndex-previousRankIndex != 1 {
				if !(i == len(h)-1 && card.Rank == Ace && previousRankIndex == 3) { // Account for A-2-3-4-5 straight
					isStraight = false
				}
			}
			previousRankIndex = currentRankIndex
		} else {
			previousRankIndex = indexOfRank(card.Rank)
		}
	}

	if isFlush && isStraight {
		score = 8000  // Straight Flush
		highCard = h[4] // Highest card
		return
	}

	if isFlush {
		score = 5000
		highCard = h[4] // Highest card in a flush is the last one assuming sorted
		return
	}

	if isStraight {
		score = 4000
		highCard = h[4] // Highest card in a straight is the last one assuming sorted
		return
	}

	// Check for Four of a Kind, Full House, Three of a Kind, Two Pair, One Pair
	var pairs, threes, fours []Rank
	for rank, count := range rankCount {
		switch count {
		case 2:
			pairs = append(pairs, rank)
		case 3:
			threes = append(threes, rank)
		case 4:
			fours = append(fours, rank)
		}
	}

	if len(fours) == 1 {
		score = 7000 // Four of a Kind
		highCard = Card{Rank: fours[0]}
		return
	}

	if len(threes) == 1 && len(pairs) == 1 {
		score = 6000 // Full House
		highCard = Card{Rank: threes[0]}
		return
	}

	if len(threes) == 1 {
		score = 3000 // Three of a Kind
		highCard = Card{Rank: threes[0]}
		return
	}

	if len(pairs) == 2 {
		score = 2000 // Two Pair
		highCard = Card{Rank: maxRank(pairs[0], pairs[1])}
		return
	}

	if len(pairs) == 1 {
		score = 1000 // One Pair
		highCard = Card{Rank: pairs[0]}
		return
	}

	// High Card
	score = 0
	highCard = h[4] // Assuming sorted, the last card is the high card
	return
}

func indexOfRank(rank Rank) int {
	for i, r := range ranks {
		if r == rank {
			return i
		}
	}
	return -1
}

func maxRank(rank1, rank2 Rank) Rank {
	if indexOfRank(rank1) > indexOfRank(rank2) {
		return rank1
	}
	return rank2
}

func payoutRatio(score int) float64 {
	switch {
	case score >= 8000:
		return 50 // Straight Flush
	case score >= 7000:
		return 25 // Four of a Kind
	case score >= 6000:
		return 9 // Full House
	case score >= 5000:
		return 6 // Flush
	case score >= 4000:
		return 4 // Straight
	case score >= 3000:
		return 3 // Three of a Kind
	case score >= 2000:
		return 2 // Two Pair
	case score >= 1000:
		return 1 // One Pair
	default:
		return 0 // High Card
	}
}

func main() {
	var balance float64 = 100
	var deck Deck
	var playerHand Hand
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Your balance: $%.2f\n", balance)
		fmt.Println("Enter bet amount:")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		bet, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}
		if bet > balance {
			fmt.Println("You cannot bet more than your balance.")
			continue
		}

		// Start the game
		balance -= bet
		deck = NewDeck()
		deck.Shuffle()
		playerHand = deck.Draw(5)

		fmt.Printf("Your hand: %s\n", playerHand)
		fmt.Println("Which cards would you like to hold? (ex: 1 3 5 for first, third, and fifth cards)")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		holds := strings.Fields(input)

		newHand := make(Hand, 0, 5)
		for _, hold := range holds {
			index, _ := strconv.Atoi(hold)
			if index < 1 || index > 5 {
				continue
			}
			newHand = append(newHand, playerHand[index-1])
		}

		// Draw new cards
		newCardsNeeded := 5 - len(newHand)
		newHand = append(newHand, deck.Draw(newCardsNeeded)...)

		// Show the final hand
		playerHand = newHand
		fmt.Printf("Your final hand: %s\n", playerHand)
		score, highCard := playerHand.Score()
		ratio := payoutRatio(score)

		if ratio > 0 {
			winnings := bet * ratio
			balance += winnings
			fmt.Printf("You win! Payout: $%.2f (High Card: %s)\n", winnings, highCard)
		} else {
			fmt.Println("You lose! Better luck next time.")
		}

		fmt.Println("Play again? (yes/no)")
		input, _ = reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(input)) != "yes" {
			break
		}
	}
}
