package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

const (
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
	Hands   []Hand
	Scores  []int
	Stati   []PlayerStatus
	Bets    []int // A bet for each hand
	IsHuman bool
}

// Game represents the state of a blackjack game.
type Game struct {
	Deck        Deck
	Player      Player
	Dealer      Player
	GameState   string // e.g., "betting", "playing", "game_over"
	PlayerBet   int
	PlayerChips int
	ActiveHand  int
}

// VisibleGame is the version of the Game struct that is sent to the client.
type VisibleGame struct {
	Player           Player   `json:"Player"`
	Dealer           Player   `json:"Dealer"`
	GameState        string   `json:"GameState"`
	PlayerChips      int      `json:"PlayerChips"`
	PlayerBet        int      `json:"PlayerBet"`
	AvailableActions []string `json:"AvailableActions"`
}

// Visible returns a version of the game state that is safe to show to the client.
func (g *Game) Visible() VisibleGame {
	visibleDealer := g.Dealer
	if g.GameState == "playing" {
		// Hide the dealer's second card
		if len(visibleDealer.Hands) > 0 && len(visibleDealer.Hands[0]) > 1 {
			visibleDealer.Hands[0] = g.Dealer.Hands[0][:1]
			visibleDealer.Scores[0] = HandScore(visibleDealer.Hands[0])
		}
	}

	return VisibleGame{
		Player:           g.Player,
		Dealer:           visibleDealer,
		GameState:        g.GameState,
		PlayerChips:      g.PlayerChips,
		PlayerBet:        g.PlayerBet,
		AvailableActions: g.getAvailableActions(),
	}
}

func (g *Game) getAvailableActions() []string {
	actions := []string{}
	if g.GameState == "betting" {
		actions = append(actions, "bet")
	}
	if g.GameState == "playing" {
		actions = append(actions, "hit", "stand")
		// Player can double down on the first two cards
		if len(g.Player.Hands[0]) == 2 {
			actions = append(actions, "doubledown")
		}
		// Player can split on a pair
		if len(g.Player.Hands) == 1 && len(g.Player.Hands[0]) == 2 && g.Player.Hands[0][0].Value == g.Player.Hands[0][1].Value {
			actions = append(actions, "split")
		}
	}
	return actions
}

// NewGame creates a new game.
func NewGame() Game {
	game := Game{
		PlayerChips: 1000, // Starting chips
		GameState:   "betting",
	}
	return game
}

// PlaceBet places a bet for the player and deals a new hand.
func (g *Game) PlaceBet(amount int) {
	if g.GameState != "betting" {
		return
	}
	if amount <= 0 || amount > g.PlayerChips {
		// Invalid bet amount
		return
	}

	g.PlayerChips -= amount
	g.PlayerBet = amount
	g.dealHand()
}

func (g *Game) dealHand() {
	deck := NewDeck()
	deck.Shuffle()

	playerHand := Hand{deck[0], deck[1]}
	dealerHand := Hand{deck[2], deck[3]}

	g.Player = Player{
		Hands:   []Hand{playerHand},
		Scores:  []int{HandScore(playerHand)},
		Stati:   []PlayerStatus{Playing},
		Bets:    []int{g.PlayerBet},
		IsHuman: true,
	}

	g.Dealer = Player{
		Hands:  []Hand{dealerHand},
		Scores: []int{HandScore(dealerHand)},
		Stati:  []PlayerStatus{Playing},
	}

	g.Deck = deck[4:]
	g.GameState = "playing"

	// Check for blackjack
	if g.Player.Scores[0] == Blackjack {
		g.Player.Stati[0] = BlackjackWin
		if g.Dealer.Scores[0] == Blackjack {
			g.Dealer.Stati[0] = BlackjackWin
			g.Player.Stati[0] = Push
			g.PlayerChips += g.PlayerBet // Push, return bet
			g.GameState = "game_over"
		} else {
			g.PlayerChips += g.PlayerBet + (g.PlayerBet*3)/2 // Blackjack pays 3:2
			g.GameState = "game_over"
		}
	} else if g.Dealer.Scores[0] == Blackjack {
		g.Dealer.Stati[0] = BlackjackWin
		g.GameState = "game_over"
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
	for score > Blackjack && aces > 0 {
		score -= 10
		aces--
	}
	return score
}

// Hit gives the current player another card.
func (g *Game) Hit() {
	if g.GameState != "playing" {
		return
	}

	player := &g.Player
	if player.Stati[g.ActiveHand] != Playing {
		return
	}

	player.Hands[g.ActiveHand] = append(player.Hands[g.ActiveHand], g.Deck[0])
	g.Deck = g.Deck[1:]
	player.Scores[g.ActiveHand] = HandScore(player.Hands[g.ActiveHand])

	if player.Scores[g.ActiveHand] > 21 {
		player.Stati[g.ActiveHand] = Bust
		g.nextHandOrDealer()
	}
}

// Stand ends the current player's turn for the current hand.
func (g *Game) Stand() {
	if g.GameState != "playing" {
		return
	}

	player := &g.Player
	if player.Stati[g.ActiveHand] != Playing {
		return
	}

	player.Stati[g.ActiveHand] = Stand
	g.nextHandOrDealer()
}

func (g *Game) nextHandOrDealer() {
	g.ActiveHand++
	if g.ActiveHand >= len(g.Player.Hands) {
		g.dealerTurn()
	}
}

// DoubleDown doubles the player's bet, deals one more card, and ends the turn.
func (g *Game) DoubleDown() {
	if g.GameState != "playing" {
		return
	}
	if g.PlayerChips < g.PlayerBet {
		// Not enough chips to double down
		return
	}

	g.PlayerChips -= g.PlayerBet
	g.Player.Bets[0] *= 2

	// Deal one more card
	player := &g.Player
	player.Hands[0] = append(player.Hands[0], g.Deck[0])
	g.Deck = g.Deck[1:]
	player.Scores[0] = HandScore(player.Hands[0])

	if player.Scores[0] > 21 {
		player.Stati[0] = Bust
		g.determineWinner()
	} else {
		player.Stati[0] = Stand
		g.dealerTurn()
	}
}

// Split splits the player's hand into two hands.
func (g *Game) Split() {
	if g.GameState != "playing" {
		return
	}
	player := &g.Player
	if len(player.Hands) != 1 || len(player.Hands[0]) != 2 || player.Hands[0][0].Value != player.Hands[0][1].Value {
		// Can only split on a pair in a single hand
		return
	}
	if g.PlayerChips < g.PlayerBet {
		// Not enough chips to split
		return
	}

	g.PlayerChips -= g.PlayerBet

	// Create two new hands
	hand1 := Hand{player.Hands[0][0], g.Deck[0]}
	hand2 := Hand{player.Hands[0][1], g.Deck[1]}
	g.Deck = g.Deck[2:]

	player.Hands = []Hand{hand1, hand2}
	player.Scores = []int{HandScore(hand1), HandScore(hand2)}
	player.Stati = []PlayerStatus{Playing, Playing}
	player.Bets = []int{g.PlayerBet, g.PlayerBet}
}

// dealerTurn plays the dealer's turn.
func (g *Game) dealerTurn() {
	// Dealer plays
	for g.Dealer.Scores[0] < DealerStand {
		g.Dealer.Hands[0] = append(g.Dealer.Hands[0], g.Deck[0])
		g.Deck = g.Deck[1:]
		g.Dealer.Scores[0] = HandScore(g.Dealer.Hands[0])
	}
	if g.Dealer.Scores[0] > 21 {
		g.Dealer.Stati[0] = Bust
	} else {
		g.Dealer.Stati[0] = Stand
	}

	// Determine the winner
	g.determineWinner()
}

// determineWinner determines the winner of the game.
func (g *Game) determineWinner() {
	dealerScore := g.Dealer.Scores[0]
	player := &g.Player

	for i := range player.Hands {
		// If player has blackjack is handled in dealHand

		// If player is bust
		if player.Stati[i] == Bust {
			player.Stati[i] = DealerWins
		} else if g.Dealer.Stati[0] == Bust {
			// If dealer is bust
			player.Stati[i] = PlayerWins
			g.PlayerChips += player.Bets[i] * 2
		} else if player.Stati[i] == Stand {
			// Compare scores
			if player.Scores[i] > dealerScore {
				player.Stati[i] = PlayerWins
				g.PlayerChips += player.Bets[i] * 2
			} else if player.Scores[i] < dealerScore {
				player.Stati[i] = DealerWins
			} else {
				player.Stati[i] = Push
				g.PlayerChips += player.Bets[i]
			}
		}
	}
	g.GameState = "game_over"
}
