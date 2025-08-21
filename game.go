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
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

func init() {
	rand.Seed(time.Now().UnixNano())
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

// VisibleGame is the version of the Game struct that is sent to the client.
type VisibleGame struct {
	Player      Hand   `json:"Player"`
	Dealer      Hand   `json:"Dealer"`
	PlayerScore int    `json:"PlayerScore"`
	DealerScore int    `json:"DealerScore"`
	State       string `json:"State"`
}

// Visible returns a version of the game state that is safe to show to the client.
func (g *Game) Visible() VisibleGame {
	if g.State != "playing" {
		return VisibleGame{
			Player:      g.Player,
			Dealer:      g.Dealer,
			PlayerScore: g.PlayerScore,
			DealerScore: g.DealerScore,
			State:       g.State,
		}
	}
	return VisibleGame{
		Player:      g.Player,
		Dealer:      g.Dealer[:1],
		PlayerScore: g.PlayerScore,
		DealerScore: HandScore(g.Dealer[:1]),
		State:       g.State,
	}
}

// NewGame creates a new game with a shuffled deck and two cards for the player and dealer.
func NewGame() Game {
	deck := NewDeck()
	deck.Shuffle()

	playerHand := Hand{deck[0], deck[2]}
	dealerHand := Hand{deck[1], deck[3]}

	playerScore := HandScore(playerHand)
	dealerScore := HandScore(dealerHand)

	state := "playing"
	if playerScore == 21 {
		if dealerScore == 21 {
			state = "push"
		} else {
			state = "player_blackjack"
		}
	} else if dealerScore == 21 {
		state = "dealer_blackjack"
	}

	return Game{
		Deck:        deck[4:],
		Player:      playerHand,
		Dealer:      dealerHand,
		PlayerScore: playerScore,
		DealerScore: dealerScore,
		State:       state,
	}
}

// HandScore calculates the score of a hand.
func HandScore(hand Hand) int {
	score := 0
	aces := 0
	for _, card := range hand {
		const (
			Ace   = 1
			Jack  = 11
			Queen = 12
			King  = 13
		)
		if card.Value >= Jack {
			score += 10
		} else if card.Value == Ace {
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

// Hit gives the player another card.
func (g *Game) Hit() {
	if g.State != "playing" {
		return
	}

	g.Player = append(g.Player, g.Deck[0])
	g.Deck = g.Deck[1:]
	g.PlayerScore = HandScore(g.Player)

	if g.PlayerScore > 21 {
		g.State = "player_busts"
	}
}

// Stand ends the player's turn and plays the dealer's turn.
func (g *Game) Stand() {
	if g.State != "playing" {
		return
	}

	// Dealer plays
	for g.DealerScore < 17 {
		g.Dealer = append(g.Dealer, g.Deck[0])
		g.Deck = g.Deck[1:]
		g.DealerScore = HandScore(g.Dealer)
	}

	// Determine the winner
	if g.DealerScore > 21 || g.PlayerScore > g.DealerScore {
		g.State = "player_wins"
	} else if g.DealerScore > g.PlayerScore {
		g.State = "dealer_wins"
	} else {
		g.State = "push"
	}
}