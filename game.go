package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

const (
	NumPlayers   = 5
	HumanPlayer  = 2
	NumCardsDeal = 2
	Blackjack    = 21
	DealerStand  = 17
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

var deterministicShuffle = false

func init() {
	if !deterministicShuffle {
		rand.Seed(time.Now().UnixNano())
	}
}

// Hand represents a player's or dealer's hand of cards.
type Hand []Card

type PlayerStatus int

const (
	Playing PlayerStatus = iota
	Bust
	Stand
	BlackjackWin
	Push
	PlayerWins
	DealerWins
)

func (s PlayerStatus) String() string {
	return [...]string{"playing", "bust", "stand", "blackjack", "push", "player_wins", "dealer_wins"}[s]
}

func (s PlayerStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// Player represents a player in the game.
type Player struct {
	Hand    Hand
	Score   int
	Status  PlayerStatus
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

	players := make([]Player, NumPlayers)
	for i := 0; i < NumPlayers; i++ {
		players[i] = Player{
			Hand:    Hand{deck[i*NumCardsDeal], deck[i*NumCardsDeal+1]},
			Status:  Playing,
			IsHuman: i == HumanPlayer, // The middle player is human
		}
		players[i].Score = HandScore(players[i].Hand)
		if players[i].Score == Blackjack {
			players[i].Status = BlackjackWin
		}
	}

	dealerHand := Hand{deck[NumPlayers*NumCardsDeal], deck[NumPlayers*NumCardsDeal+1]}
	dealer := Player{
		Hand:   dealerHand,
		Score:  HandScore(dealerHand),
		Status: Playing,
	}
	if dealer.Score == Blackjack {
		dealer.Status = BlackjackWin
	}

	game := Game{
		Deck:    deck[NumPlayers*NumCardsDeal+NumCardsDeal:],
		Players: players,
		Dealer:  dealer,
		State:   "playing",
		Turn:    -1, // Start before the first player
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
	for score > Blackjack && aces > 0 {
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
	if player.Status != Playing {
		return
	}

	player.Hand = append(player.Hand, g.Deck[0])
	g.Deck = g.Deck[1:]
	player.Score = HandScore(player.Hand)

	if player.Score > 21 {
		player.Status = Bust
	}
}

// Stand ends the current player's turn.
func (g *Game) Stand() {
	if g.State != "playing" {
		return
	}

	player := &g.Players[g.Turn]
	if player.Status != Playing {
		return
	}

	player.Status = Stand
}

// NextTurn moves to the next player or the dealer's turn.
func (g *Game) NextTurn() {
	g.Turn++

	for g.Turn < len(g.Players) {
		player := &g.Players[g.Turn]
		if player.IsHuman && player.Status == Playing {
			// It's the human player's turn
			return
		}

		if player.Status == Playing {
			// It's a computer player's turn
			for player.Score < DealerStand {
				player.Hand = append(player.Hand, g.Deck[0])
				g.Deck = g.Deck[1:]
				player.Score = HandScore(player.Hand)
			}
			if player.Score > 21 {
				player.Status = Bust
			} else {
				player.Status = Stand
			}
		}
		g.Turn++
	}

	if g.Turn >= len(g.Players) {
		g.dealerTurn()
	}
}

// dealerTurn plays the dealer's turn.
func (g *Game) dealerTurn() {
	// Dealer plays
	for g.Dealer.Score < DealerStand {
		g.Dealer.Hand = append(g.Dealer.Hand, g.Deck[0])
		g.Deck = g.Deck[1:]
		g.Dealer.Score = HandScore(g.Dealer.Hand)
	}
	if g.Dealer.Score > 21 {
		g.Dealer.Status = Bust
	} else {
		g.Dealer.Status = Stand
	}


	// Determine the winner
	g.determineWinner()
}

// determineWinner determines the winner of the game.
func (g *Game) determineWinner() {
	dealerScore := g.Dealer.Score
	for i := range g.Players {
		player := &g.Players[i]

		// If player has blackjack
		if player.Status == BlackjackWin {
			if g.Dealer.Status == BlackjackWin {
				player.Status = Push // Both have blackjack
			} else {
				player.Status = PlayerWins // Player has blackjack, dealer doesn't
			}
			continue
		}

		// If dealer has blackjack
		if g.Dealer.Status == BlackjackWin {
			if player.Status != Bust {
				player.Status = DealerWins
			}
			continue
		}

		// If player is bust
		if player.Status == Bust {
			player.Status = DealerWins
			continue
		}

		// If dealer is bust
		if g.Dealer.Status == Bust {
			player.Status = PlayerWins
			continue
		}

		// Compare scores
		if player.Status == Stand {
			if player.Score > dealerScore {
				player.Status = PlayerWins
			} else if player.Score < dealerScore {
				player.Status = DealerWins
			} else {
				player.Status = Push
			}
		}
	}
	g.State = "game_over"
}
