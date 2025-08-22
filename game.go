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

// Player represents a player in the game.
type Player struct {
	Hand    Hand
	Score   int
	Status  string // e.g., "playing", "bust", "stand", "blackjack"
	IsHuman bool
}

// Game represents the state of a blackjack game.
type Game struct {
	Deck    Deck
	Players []Player
	Dealer  Player
	State   string // e.g., "playing", "game_over"
	Turn    int    // Index of the current player in the Players slice
}

// VisibleGame is the version of the Game struct that is sent to the client.
type VisibleGame struct {
	Players []Player `json:"Players"`
	Dealer  Player   `json:"Dealer"`
	State   string   `json:"State"`
	Turn    int      `json:"Turn"`
}

// Visible returns a version of the game state that is safe to show to the client.
func (g *Game) Visible() VisibleGame {
	if g.State != "playing" {
		return VisibleGame{
			Players: g.Players,
			Dealer:  g.Dealer,
			State:   g.State,
			Turn:    g.Turn,
		}
	}

	// Hide the dealer's second card
	visibleDealer := g.Dealer
	visibleDealer.Hand = g.Dealer.Hand[:1]
	visibleDealer.Score = HandScore(visibleDealer.Hand)

	return VisibleGame{
		Players: g.Players,
		Dealer:  visibleDealer,
		State:   g.State,
		Turn:    g.Turn,
	}
}

// NewGame creates a new game with a shuffled deck and two cards for each player and the dealer.
func NewGame() Game {
	deck := NewDeck()
	deck.Shuffle()

	players := make([]Player, 5)
	for i := 0; i < 5; i++ {
		players[i] = Player{
			Hand:    Hand{deck[i*2], deck[i*2+1]},
			Status:  "playing",
			IsHuman: i == 2, // The middle player is human
		}
		players[i].Score = HandScore(players[i].Hand)
		if players[i].Score == 21 {
			players[i].Status = "blackjack"
		}
	}

	dealerHand := Hand{deck[10], deck[11]}
	dealer := Player{
		Hand:   dealerHand,
		Score:  HandScore(dealerHand),
		Status: "playing",
	}
	if dealer.Score == 21 {
		dealer.Status = "blackjack"
	}

	game := Game{
		Deck:    deck[12:],
		Players: players,
		Dealer:  dealer,
		State:   "playing",
		Turn:    0, // Start with the first player
	}

	return game
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

// Hit gives the current player another card.
func (g *Game) Hit() {
	if g.State != "playing" {
		return
	}

	player := &g.Players[g.Turn]
	if player.Status != "playing" {
		return
	}

	player.Hand = append(player.Hand, g.Deck[0])
	g.Deck = g.Deck[1:]
	player.Score = HandScore(player.Hand)

	if player.Score > 21 {
		player.Status = "bust"
		g.NextTurn()
	}
}

// Stand ends the current player's turn.
func (g *Game) Stand() {
	if g.State != "playing" {
		return
	}

	player := &g.Players[g.Turn]
	if player.Status != "playing" {
		return
	}

	player.Status = "stand"
	g.NextTurn()
}

// NextTurn moves to the next player or the dealer's turn.
func (g *Game) NextTurn() {
	for g.Turn < len(g.Players) && g.Players[g.Turn].Status != "playing" {
		g.Turn++
	}

	if g.Turn >= len(g.Players) {
		g.dealerTurn()
	}
}

// dealerTurn plays the dealer's turn.
func (g *Game) dealerTurn() {
	// Dealer plays
	for g.Dealer.Score < 17 {
		g.Dealer.Hand = append(g.Dealer.Hand, g.Deck[0])
		g.Deck = g.Deck[1:]
		g.Dealer.Score = HandScore(g.Dealer.Hand)
	}

	// Determine the winner
	g.determineWinner()
}

// determineWinner determines the winner of the game.
func (g *Game) determineWinner() {
	dealerScore := g.Dealer.Score
	for i := range g.Players {
		player := &g.Players[i]
		if player.Status == "blackjack" {
			if g.Dealer.Status == "blackjack" {
				player.Status = "push"
			} else {
				player.Status = "player_wins"
			}
		} else if player.Status == "bust" {
			player.Status = "dealer_wins"
		} else if player.Status == "playing" || player.Status == "stand" {
			if dealerScore > 21 || player.Score > dealerScore {
				player.Status = "player_wins"
			} else if dealerScore > player.Score {
				player.Status = "dealer_wins"
			} else {
				player.Status = "push"
			}
		}
	}
	g.State = "game_over"
}
