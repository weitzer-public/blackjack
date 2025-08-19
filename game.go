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
	i := 0
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
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

// Hand represents a player's or dealer's hand of cards.
type Hand []Card

// Game represents the state of a blackjack game.
type Game struct {
	Deck   Deck
	Player Hand
	Dealer Hand
}

// NewGame creates a new game with a shuffled deck and two cards for the player and dealer.
func NewGame() Game {
	deck := NewDeck()
	deck.Shuffle()

	playerHand := Hand{deck[0], deck[2]}
	dealerHand := Hand{deck[1], deck[3]}

	game := Game{
		Deck:   deck[4:],
		Player: playerHand,
		Dealer: dealerHand,
	}

	return game
}

// HandScore calculates the score of a hand.
func HandScore(hand Hand) int {
	score := 0
	aces := 0
	for _, card := range hand {
		if card.Value > 10 {
			score += 10
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
