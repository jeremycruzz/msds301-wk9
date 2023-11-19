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
	numSuits = 4
	numRanks = 13
	numCards = numSuits * numRanks
	handSize = 5
)

var suits = [numSuits]string{"Hearts", "Diamonds", "Clubs", "Spades"}
var ranks = [numRanks]string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}

type Card struct {
	rank string
	suit string
}

type Deck []Card

func NewDeck() Deck {
	deck := make(Deck, numCards)
	i := 0
	for _, suit := range suits {
		for _, rank := range ranks {
			deck[i] = Card{rank, suit}
			i++
		}
	}
	return deck
}

func (deck Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

func (deck *Deck) Deal(n int) []Card {
	hand := make([]Card, n)
	copy(hand, (*deck)[:n])
	*deck = (*deck)[n:]
	return hand
}

func (deck *Deck) ReturnCards(cards []Card) {
	*deck = append(*deck, cards...)
}

type Player struct {
	hand  []Card
	money int
	bet   int
}

func NewPlayer(money int) Player {
	return Player{
		hand:  []Card{},
		money: money,
	}
}

func (p *Player) ReceiveCards(cards []Card) {
	p.hand = cards
}

func (p *Player) ShowHand() {
	handStr := "Your hand: "
	for _, card := range p.hand {
		handStr += fmt.Sprintf("[%s of %s] ", card.rank, card.suit)
	}
	fmt.Println(handStr)
}

func (p *Player) PlaceBet(betAmount int) {
	p.bet = betAmount
	p.money -= betAmount
}

func (p Player) HasMoney() bool {
	return p.money > 0
}

func evaluateHand(hand []Card) (string, int) {
	rankOccurrences := make(map[string]int)
	suitOccurrences := make(map[string]int)
	for _, card := range hand {
		rankOccurrences[card.rank]++
		suitOccurrences[card.suit]++
	}

	isFlush := len(suitOccurrences) == 1
	isStraight := false
	straightStartIndex := -1

	for i := 0; i <= numRanks-handSize; i++ {
		count := 0
		for j := 0; j < handSize; j++ {
			if rankOccurrences[ranks[i+j]] > 0 {
				count++
			}
		}
		if count == handSize {
			isStraight = true
			straightStartIndex = i
			break
		}
	}

	if isStraight && isFlush && straightStartIndex == numRanks-handSize {
		return "Royal Flush", 9
	}

	if isStraight && isFlush {
		return "Straight Flush", 8
	}

	highestRankCount := 0
	for _, count := range rankOccurrences {
		if count > highestRankCount {
			highestRankCount = count
		}
	}

	switch highestRankCount {
	case 4:
		return "Four of a Kind", 7
	case 3:
		if len(rankOccurrences) == 2 {
			return "Full House", 6
		}
		return "Three of a Kind", 3
	case 2:
		if len(rankOccurrences) == 3 {
			return "Two Pair", 2
		}
		return "One Pair", 1
	default:
		if isFlush {
			return "Flush", 5
		}
		if isStraight {
			return "Straight", 4
		}
	}

	return "High Card", 0
}

func (p *Player) EvaluateHand() (string, int) {
	sort.Slice(p.hand, func(i int, j int) bool {
		return cardValue(p.hand[i]) < cardValue(p.hand[j])
	})
	handType, rank := evaluateHand(p.hand)
	fmt.Printf("You have a %s\n", handType)
	return handType, rank
}

func cardValue(card Card) int {
	valueMap := map[string]int{
		"2": 2, "3": 3, "4": 4, "5": 5, "6": 6,
		"7": 7, "8": 8, "9": 9, "T": 10,
		"J": 11, "Q": 12, "K": 13, "A": 14,
	}
	return valueMap[card.rank] + numRanks*indexOf(suits[:], card.suit)
}

func indexOf(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

func (p *Player) Draw(deck *Deck, discardIndices ...int) {
	cardsToDiscard := make([]Card, len(discardIndices))
	for i, index := range discardIndices {
		cardsToDiscard[i] = p.hand[index]
	}

	p.hand = append(p.hand[:0], p.hand[:len(p.hand)]...)
	deck.ReturnCards(cardsToDiscard)
	deck.Shuffle()
	newCards := deck.Deal(len(cardsToDiscard))
	p.hand = append(p.hand, newCards...)

	p.EvaluateHand()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to GoPoker!")
	fmt.Print("Enter your starting money: $")
	input, _ := reader.ReadString('\n')
	startMoney, err := strconv.Atoi(strings.TrimSpace(input))
	for err != nil || startMoney <= 0 {
		fmt.Println("Invalid amount, please enter a positive number for your money.")
		input, _ = reader.ReadString('\n')
		startMoney, err = strconv.Atoi(strings.TrimSpace(input))
	}

	player := NewPlayer(startMoney)
	deck := NewDeck()
	deck.Shuffle()

	keepPlaying := true
	for keepPlaying && player.HasMoney() {
		fmt.Printf("You have $%d. How much would you like to bet? $", player.money)
		input, _ := reader.ReadString('\n')
		bet, _ := strconv.Atoi(strings.TrimSpace(input))
		player.PlaceBet(bet)

		deck.Shuffle()
		player.ReceiveCards(deck.Deal(handSize))
		player.ShowHand()
		player.EvaluateHand()

		fmt.Print("Enter card positions to discard (e.g., 1 3 5 or 'none' to keep): ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "none" {
			fmt.Println("No cards discarded.")
		} else {
			discardIndexes := parseDiscardIndexes(input)
			player.Draw(&deck, discardIndexes...)
		}

		handType, _ := player.EvaluateHand()

		fmt.Printf("Final hand: %s ($%d bet)\n", handType, player.bet)

		fmt.Printf("Would you like to play another round? (yes/no): ")
		input, _ = reader.ReadString('\n')
		if strings.TrimSpace(input) != "yes" {
			keepPlaying = false
		}

		deck.ReturnCards(player.hand)
		player.hand = []Card{}
	}
	if player.money <= 0 {
		fmt.Println("You've run out of money. Game over.")
	} else {
		fmt.Println("Thank you for playing GoPoker!")
	}
	fmt.Printf("You ended with $%d\n", player.money)
}

func parseDiscardIndexes(input string) []int {
	var indexes []int
	if input != "none" {
		inputs := strings.Fields(input)
		for _, s := range inputs {
			index, _ := strconv.Atoi(s)
			if index >= 1 && index <= handSize {
				indexes = append(indexes, index-1)
			}
		}
	}
	return indexes
}
