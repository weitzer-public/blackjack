package main

import (
	"encoding/json"
	"math/rand"
	"strconv"
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

type Player struct {
	Name    string
	Hands   []Hand
	Scores  []int
	Stati   []PlayerStatus
	Bets    []int // A bet for each hand
	IsHuman bool
	Chips   int
}

// Game represents the state of a blackjack game.
type Game struct {
	Deck         Deck
	Players      []Player
	Dealer       Player
	GameState    string // e.g., "betting", "playing", "game_over"
	HumanPlayer  *Player
	ActivePlayer int
	ActiveHand   int
}

// VisibleGame is the version of the Game struct that is sent to the client.
type VisibleGame struct {
	Players          []Player `json:"Players"`
	Dealer           Player   `json:"Dealer"`
	GameState        string   `json:"GameState"`
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
		Players:          g.Players,
		Dealer:           visibleDealer,
		GameState:        g.GameState,
		AvailableActions: g.getAvailableActions(),
	}
}

func (g *Game) getAvailableActions() []string {
	actions := []string{}
	if g.GameState == "betting" {
		actions = append(actions, "bet")
	}
	if g.GameState == "playing" && g.HumanPlayer != nil && g.Players[g.ActivePlayer].IsHuman {
		actions = append(actions, "hit", "stand")
		// Player can double down on the first two cards
		if len(g.HumanPlayer.Hands[g.ActiveHand]) == 2 {
			actions = append(actions, "doubledown")
		}
		// Player can split on a pair
		if len(g.HumanPlayer.Hands) == 1 && len(g.HumanPlayer.Hands[0]) == 2 && g.HumanPlayer.Hands[0][0].Value == g.HumanPlayer.Hands[0][1].Value {
			actions = append(actions, "split")
		}
	}
	return actions
}

func NewGame() Game {
	game := Game{
		GameState: "betting",
		Players:   make([]Player, 5),
	}

	// Initialize players
	for i := 0; i < 5; i++ {
		game.Players[i] = Player{
			Name:    "Player " + strconv.Itoa(i+1),
			Chips:   1000,
			IsHuman: i == 0, // First player is human
		}
	}
	game.HumanPlayer = &game.Players[0]

	return game
}

// PlaceBet places a bet for the player and deals a new hand.
func (g *Game) PlaceBet(amount int) {
	if g.GameState != "betting" {
		return
	}
	if amount <= 0 || amount > g.HumanPlayer.Chips {
		// Invalid bet amount
		return
	}

	g.HumanPlayer.Chips -= amount
	g.HumanPlayer.Bets = []int{amount}

	// AI players place a fixed bet
	for i := range g.Players {
		if !g.Players[i].IsHuman {
			g.Players[i].Bets = []int{10}
			g.Players[i].Chips -= 10
		}
	}

	g.dealHand()
}

func (g *Game) dealHand() {
	deck := NewDeck()
	deck.Shuffle()

	// Deal to players
	for i := range g.Players {
		hand := Hand{deck[0], deck[1]}
		g.Players[i].Hands = []Hand{hand}
		g.Players[i].Scores = []int{HandScore(hand)}
		g.Players[i].Stati = []PlayerStatus{Playing}
		deck = deck[2:]
	}

	// Deal to dealer
	dealerHand := Hand{deck[0], deck[1]}
	g.Dealer = Player{
		Hands:  []Hand{dealerHand},
		Scores: []int{HandScore(dealerHand)},
		Stati:  []PlayerStatus{Playing},
	}
	deck = deck[2:]

	g.Deck = deck
	g.GameState = "playing"
	g.ActivePlayer = 0

	// Check for blackjack
	for i := range g.Players {
		if g.Players[i].Scores[0] == Blackjack {
			g.Players[i].Stati[0] = BlackjackWin
		}
	}

	if g.Dealer.Scores[0] == Blackjack {
		g.Dealer.Stati[0] = BlackjackWin
	}

	// If dealer or all players have blackjack, game is over
	gameOver := true
	for i := range g.Players {
		if g.Players[i].Stati[0] != BlackjackWin {
			gameOver = false
			break
		}
	}

	if g.Dealer.Stati[0] == BlackjackWin || gameOver {
		g.determineWinner()
	} else {
		g.nextPlayer()
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

func (g *Game) nextPlayer() {
	if g.ActivePlayer >= len(g.Players) {
		g.dealerTurn()
		return
	}

	player := &g.Players[g.ActivePlayer]
	if player.Stati[0] == BlackjackWin {
		g.ActivePlayer++
		g.nextPlayer()
		return
	}

	if !player.IsHuman {
		// AI's turn
		g.playAITurn()
		g.ActivePlayer++
		g.nextPlayer()
	}
	// It's human's turn, wait for input
}

func (g *Game) playAITurn() {
	player := &g.Players[g.ActivePlayer]
	for HandScore(player.Hands[0]) < 17 {
		player.Hands[0] = append(player.Hands[0], g.Deck[0])
		g.Deck = g.Deck[1:]
		player.Scores[0] = HandScore(player.Hands[0])
	}
	if HandScore(player.Hands[0]) > 21 {
		player.Stati[0] = Bust
	} else {
		player.Stati[0] = Stand
	}
}

// Hit gives the current player another card.
func (g *Game) Hit() {
	if g.GameState != "playing" {
		return
	}

	player := &g.Players[g.ActivePlayer]
	if !player.IsHuman || player.Stati[g.ActiveHand] != Playing {
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

	player := &g.Players[g.ActivePlayer]
	if !player.IsHuman || player.Stati[g.ActiveHand] != Playing {
		return
	}

	player.Stati[g.ActiveHand] = Stand
	g.nextHandOrDealer()
}

func (g *Game) nextHandOrDealer() {
	g.ActiveHand++
	if g.ActiveHand >= len(g.Players[g.ActivePlayer].Hands) {
		g.ActivePlayer++
		g.ActiveHand = 0
		g.nextPlayer()
	}
}

// DoubleDown doubles the player's bet, deals one more card, and ends the turn.
func (g *Game) DoubleDown() {
	if g.GameState != "playing" {
		return
	}
	player := &g.Players[g.ActivePlayer]
	if !player.IsHuman {
		return
	}
	if player.Chips < player.Bets[0] {
		// Not enough chips to double down
		return
	}

	player.Chips -= player.Bets[0]
	player.Bets[0] *= 2

	// Deal one more card
	player.Hands[0] = append(player.Hands[0], g.Deck[0])
	g.Deck = g.Deck[1:]
	player.Scores[0] = HandScore(player.Hands[0])

	if player.Scores[0] > 21 {
		player.Stati[0] = Bust
	} else {
		player.Stati[0] = Stand
	}
	g.ActivePlayer++
	g.nextPlayer()
}

// Split splits the player's hand into two hands.
func (g *Game) Split() {
	if g.GameState != "playing" {
		return
	}
	player := &g.Players[g.ActivePlayer]
	if !player.IsHuman {
		return
	}
	if len(player.Hands) != 1 || len(player.Hands[0]) != 2 || player.Hands[0][0].Value != player.Hands[0][1].Value {
		// Can only split on a pair in a single hand
		return
	}
	if player.Chips < player.Bets[0] {
		// Not enough chips to split
		return
	}

	player.Chips -= player.Bets[0]

	// Create two new hands
	hand1 := Hand{player.Hands[0][0], g.Deck[0]}
	hand2 := Hand{player.Hands[0][1], g.Deck[1]}
	g.Deck = g.Deck[2:]

	player.Hands = []Hand{hand1, hand2}
	player.Scores = []int{HandScore(hand1), HandScore(hand2)}
	player.Stati = []PlayerStatus{Playing, Playing}
	player.Bets = []int{player.Bets[0], player.Bets[0]}
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

	for i := range g.Players {
		player := &g.Players[i]
		for j := range player.Hands {
			// If player has blackjack
			if player.Stati[j] == BlackjackWin {
				if g.Dealer.Stati[0] == BlackjackWin {
					player.Stati[j] = Push
					player.Chips += player.Bets[j]
				} else {
					player.Chips += player.Bets[j] + (player.Bets[j]*3)/2
				}
			} else if player.Stati[j] == Bust {
				player.Stati[j] = DealerWins
			} else if g.Dealer.Stati[0] == Bust {
				// If dealer is bust
				player.Stati[j] = PlayerWins
				player.Chips += player.Bets[j] * 2
			} else if player.Stati[j] == Stand {
				// Compare scores
				if player.Scores[j] > dealerScore {
					player.Stati[j] = PlayerWins
					player.Chips += player.Bets[j] * 2
				} else if player.Scores[j] < dealerScore {
					player.Stati[j] = DealerWins
				} else {
					player.Stati[j] = Push
					player.Chips += player.Bets[j]
				}
			}
		}
	}
	g.GameState = "game_over"
}