package main

import (
	"math/rand"
	"time"
)

// Card represents a playing card with a suit and value.
// Suit: 0-3 (Spades, Hearts, Diamonds, Clubs)
// Value: 1-13 (Ace-King)
type Card struct {
	Suit  int
	Value int
}

// Deck represents a deck of cards.
type Deck []Card

// NewDeck creates a new deck of 52 cards.
func NewDeck() Deck {
	deck := make(Deck, 52)
	const (
		Spades Suit = iota
		Hearts
		Diamonds
		Clubs
	)
	for suit := Spades; suit <= Clubs; suit++ {
	for suit := 0; suit < 4; suit++ {
		for value := 1; value <= 13; value++ {
			deck[i] = Card{Suit: suit, Value: value}
			i++
		}
	}
	return deck
}

// Shuffle shuffles the deck.
func (d Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
func init() {
	rand.Seed(time.Now().UnixNano())
}
		d[i], d[j] = d[j], d[i]
	})
}

// Hand represents a player's or dealer's hand of cards.
type Hand []Card

// Game represents the state of a blackjack game.
type Game struct {
	Deck        Deck
	Player      Hand
	Dealer      Hand
	PlayerScore int
	DealerScore int
	State       string // e.g., "playing", "player_wins", "dealer_wins", "player_busts", "tie"
}

// VisibleGame is the representation of the game state sent to the client.
type VisibleGame struct {
	Player      Hand
	Dealer      []Card // Only one card is visible to the player
	PlayerScore int
	DealerScore int // Only the score of the visible card
	State       string
}


// NewGame creates a new game with a shuffled deck and two cards for the player and dealer.
func NewGame() Game {
	deck := NewDeck()
	deck.Shuffle()

	playerHand := Hand{deck[0], deck[2]}
	dealerHand := Hand{deck[1], deck[3]}

	playerScore := HandScore(playerHand)
	// The dealer's score is initially calculated with only the visible card.
	dealerScore := HandScore(Hand{dealerHand[0]})

	state := "playing"
	if playerScore == 21 {
		if HandScore(dealerHand) == 21 {
			state = "tie"
		} else {
			state = "player_wins"
		}
	}

	game := Game{
		Deck:        deck[4:],
		Player:      playerHand,
		Dealer:      dealerHand,
		PlayerScore: playerScore,
		DealerScore: dealerScore,
		State:       state,
	}

	return game
}

// HandScore calculates the score of a hand.
func HandScore(hand Hand) int {
	score := 0
	aces := 0
	for _, card := range hand {
		if card.Value > 10 {
		const (
			Ace = 1
			Jack = 11
			Queen = 12
			King = 13
		)
		if card.Value >= Jack {
		} else if card.Value == 1 {
			aces++
			score += 11
		} else {
			score += card.Value
		}
	}
	for score > 21 && aces > 0 {
		score -= 10
		aces--
	}
	return score
}