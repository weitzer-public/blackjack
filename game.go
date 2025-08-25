
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
	NumAIPlayers = 4
	AIPlayerBet  = 10
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
	Name    string
	Hands   []Hand
	Scores  []int
	Stati   []PlayerStatus
	Bets    []int // A bet for each hand
	IsAI    bool
	IsTurn  bool
	Result  string // e.g., "Bust", "Blackjack", "Win", "Loss", "Push"
	Chips   int
}

// Game represents the state of a blackjack game.
type Game struct {
	Deck               Deck
	Players            []Player
	Dealer             Player
	GameState          string // e.g., "betting", "playing", "game_over"
	PlayerChips        int
	ActiveHand         int
	CurrentPlayerIndex int
}

// VisibleGame is the version of the Game struct that is sent to the client.
type VisibleGame struct {
	Players          []Player `json:"Players"`
	Dealer           Player   `json:"Dealer"`
	GameState        string   `json:"GameState"`
	PlayerChips      int      `json:"PlayerChips"`
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
		PlayerChips:      g.PlayerChips,
		AvailableActions: g.getAvailableActions(),
	}
}

func (g *Game) getAvailableActions() []string {
	actions := []string{}
	if g.GameState == "betting" {
		actions = append(actions, "bet")
	}
	if g.GameState == "playing" && g.CurrentPlayerIndex < len(g.Players) && !g.Players[g.CurrentPlayerIndex].IsAI {
		player := g.Players[g.CurrentPlayerIndex]
		actions = append(actions, "hit", "stand")
		// Player can double down on the first two cards
		if len(player.Hands[0]) == 2 {
			actions = append(actions, "doubledown")
		}
		// Player can split on a pair
		if len(player.Hands) == 1 && len(player.Hands[0]) == 2 && player.Hands[0][0].Value == player.Hands[0][1].Value {
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
	game.Players = make([]Player, NumAIPlayers+1)
	game.Players[0] = Player{Name: "Player", IsAI: false, Chips: game.PlayerChips}
	for i := 1; i <= NumAIPlayers; i++ {
		game.Players[i] = Player{Name: "AI " + strconv.Itoa(i), IsAI: true, Chips: 1000}
	}
	return game
}

// PlaceBet places a bet for the player and deals a new hand.
func (g *Game) PlaceBet(amount int) {
	if g.GameState != "betting" {
		return
	}
	if amount <= 0 || amount > g.Players[0].Chips {
		// Invalid bet amount
		return
	}

	g.Players[0].Chips -= amount
	g.Players[0].Bets = []int{amount}

	// AI players place their bets
	for i := 1; i < len(g.Players); i++ {
		g.Players[i].Bets = []int{AIPlayerBet}
		g.Players[i].Chips -= AIPlayerBet
	}

	g.dealHand()
}

func (g *Game) dealHand() {
	deck := NewDeck()
	deck.Shuffle()

	// Deal to players
	for i := range g.Players {
		hand := Hand{deck[i*2], deck[i*2+1]}
		g.Players[i].Hands = []Hand{hand}
		g.Players[i].Scores = []int{HandScore(hand)}
		g.Players[i].Stati = []PlayerStatus{Playing}
	}

	// Deal to dealer
	dealerHand := Hand{deck[len(g.Players)*2], deck[len(g.Players)*2+1]}
	g.Dealer = Player{
		Name:   "Dealer",
		Hands:  []Hand{dealerHand},
		Scores: []int{HandScore(dealerHand)},
		Stati:  []PlayerStatus{Playing},
	}

	g.Deck = deck[len(g.Players)*2+2:]
	g.GameState = "playing"
	g.CurrentPlayerIndex = 0

	// Check for blackjacks
	for i := range g.Players {
		if g.Players[i].Scores[0] == Blackjack {
			g.Players[i].Stati[0] = BlackjackWin
		}
	}
	if g.Dealer.Scores[0] == Blackjack {
		g.Dealer.Stati[0] = BlackjackWin
	}

	// If human player has blackjack, their turn is over
	if g.Players[0].Stati[0] == BlackjackWin {
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

// Hit gives the current player another card.
func (g *Game) Hit() {
	if g.GameState != "playing" {
		return
	}

	player := &g.Players[g.CurrentPlayerIndex]
	if player.Stati[g.ActiveHand] != Playing {
		return
	}

	player.Hands[g.ActiveHand] = append(player.Hands[g.ActiveHand], g.Deck[0])
	g.Deck = g.Deck[1:]
	player.Scores[g.ActiveHand] = HandScore(player.Hands[g.ActiveHand])

	if player.Scores[g.ActiveHand] > 21 {
		player.Stati[g.ActiveHand] = Bust
		g.nextPlayer()
	}
}

// Stand ends the current player's turn for the current hand.
func (g *Game) Stand() {
	if g.GameState != "playing" {
		return
	}

	player := &g.Players[g.CurrentPlayerIndex]
	if player.Stati[g.ActiveHand] != Playing {
		return
	}

	player.Stati[g.ActiveHand] = Stand
	g.nextPlayer()
}

func (g *Game) nextPlayer() {
	g.ActiveHand = 0
	g.CurrentPlayerIndex++
	if g.CurrentPlayerIndex >= len(g.Players) {
		g.dealerTurn()
	} else {
		if g.Players[g.CurrentPlayerIndex].IsAI {
			g.aiTurn()
		}
	}
}

func (g *Game) aiTurn() {
	player := &g.Players[g.CurrentPlayerIndex]
	for player.Scores[0] < DealerStand {
		g.Hit()
	}
	if player.Stati[0] == Playing {
		g.Stand()
	}
}

// DoubleDown doubles the player's bet, deals one more card, and ends the turn.
func (g *Game) DoubleDown() {
	if g.GameState != "playing" {
		return
	}
	player := &g.Players[g.CurrentPlayerIndex]
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
	g.nextPlayer()
}

// Split splits the player's hand into two hands.
func (g *Game) Split() {
	if g.GameState != "playing" {
		return
	}
	player := &g.Players[g.CurrentPlayerIndex]
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
				player.Stati[j] = PlayerWins
				player.Chips += player.Bets[j] * 2
			} else if player.Stati[j] == Stand {
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
	g.PlayerChips = g.Players[0].Chips
	g.GameState = "game_over"
}
