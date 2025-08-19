package main

import (
	"math/rand"
	"time"
)

// Card represents a playing card with a suit and value.
type Card struct {
	Suit  string
	Value string
	Rank  int
}

// Deck represents a deck of cards.
type Deck []Card

// NewDeck creates a new deck of 52 cards.
func NewDeck() Deck {
	deck := make(Deck, 52)
	suits := []string{"Spades", "Hearts", "Diamonds", "Clubs"}
	values := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	ranks := []int{11, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10}

	i := 0
	for _, suit := range suits {
		for j, value := range values {
			deck[i] = Card{Suit: suit, Value: value, Rank: ranks[j]}
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
	dealerScore := HandScore(dealerHand)

	state := "playing"
	if playerScore == 21 {
		if dealerScore == 21 {
			state = "tie"
		} else {
			state = "player_wins"
		}
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
		if card.Rank == 11 {
			aces++
		}
		score += card.Rank
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

	// Determine winner
	if g.DealerScore > 21 || g.PlayerScore > g.DealerScore {
		g.State = "player_wins"
	} else if g.DealerScore > g.PlayerScore {
		g.State = "dealer_wins"
	} else {
		g.State = "tie"
	}
}

// ToVisible converts a Game to a VisibleGame for the client.
func (g *Game) ToVisible() VisibleGame {
	if g.State == "playing" {
		return VisibleGame{
			Player:      g.Player,
			Dealer:      []Card{g.Dealer[0]}, // Only show one card
			PlayerScore: g.PlayerScore,
			DealerScore: HandScore(Hand{g.Dealer[0]}),
			State:       g.State,
		}
	}

	return VisibleGame{
		Player:      g.Player,
		Dealer:      g.Dealer,
		PlayerScore: g.PlayerScore,
		DealerScore: g.DealerScore,
		State:       g.State,
	}
}
