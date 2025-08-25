package main

import (
	"math/rand"
	"testing"
)

func init() {
	deterministicShuffle = true
	rand.Seed(0)
}

func TestHandScore(t *testing.T) {
	testCases := []struct {
		hand Hand
		want int
	}{
		{hand: Hand{{Value: 10}, {Value: 1}}, want: 21},
		{hand: Hand{{Value: 5}, {Value: 5}, {Value: 10}}, want: 20},
		{hand: Hand{{Value: 1}, {Value: 1}, {Value: 1}}, want: 13},
		{hand: Hand{{Value: 12}, {Value: 13}, {Value: 1}}, want: 21},
		{hand: Hand{{Value: 1}, {Value: 10}, {Value: 10}}, want: 21},
	}

	for _, tc := range testCases {
		got := HandScore(tc.hand)
		if got != tc.want {
			t.Errorf("HandScore(%v) = %d; want %d", tc.hand, got, tc.want)
		}
	}
}

func TestNewGame(t *testing.T) {
	game := NewGame()

	if len(game.Players) != 5 {
		t.Errorf("Expected 5 players, but got %d", len(game.Players))
	}

	if game.Players[0].Chips != 1000 {
		t.Errorf("Expected player to have 1000 chips, but got %d", game.Players[0].Chips)
	}

	if game.GameState != "betting" {
		t.Errorf("Expected game state to be 'betting', but got %s", game.GameState)
	}
}

func TestPlaceBet(t *testing.T) {
	game := NewGame()
	game.PlaceBet(100)

	if game.HumanPlayer.Chips != 900 {
		t.Errorf("Expected player to have 900 chips, but got %d", game.HumanPlayer.Chips)
	}

	if game.HumanPlayer.Bets[0] != 100 {
		t.Errorf("Expected player bet to be 100, but got %d", game.HumanPlayer.Bets[0])
	}

	if game.GameState != "playing" {
		t.Errorf("Expected game state to be 'playing', but got %s", game.GameState)
	}

	if len(game.HumanPlayer.Hands[0]) != 2 {
		t.Errorf("Expected player to have 2 cards, but got %d", len(game.HumanPlayer.Hands[0]))
	}

	if len(game.Dealer.Hands[0]) != 2 {
		t.Errorf("Expected dealer to have 2 cards, but got %d", len(game.Dealer.Hands[0]))
	}
}

func TestHit(t *testing.T) {
	game := NewGame()
	game.PlaceBet(100)
	game.Hit()

	if len(game.HumanPlayer.Hands[0]) != 3 {
		t.Errorf("Expected player to have 3 cards, but got %d", len(game.HumanPlayer.Hands[0]))
	}
}

func TestStand(t *testing.T) {
	game := NewGame()
	game.PlaceBet(100)
	game.Stand()

	if game.GameState != "game_over" {
		t.Errorf("Expected game state to be 'game_over', but got %s", game.GameState)
	}
}

func TestAITurn(t *testing.T) {
	game := NewGame()
	game.PlaceBet(100)

	// Find an AI player
	var aiPlayer *Player
	for i := range game.Players {
		if !game.Players[i].IsHuman {
			aiPlayer = &game.Players[i]
			break
		}
	}

	// Force AI player to have a low score
	aiPlayer.Hands[0] = Hand{{Value: 2}, {Value: 3}}
	aiPlayer.Scores[0] = HandScore(aiPlayer.Hands[0])

	game.ActivePlayer = 1 // Set active player to the AI
	game.playAITurn()

	if aiPlayer.Stati[0] != Stand && aiPlayer.Stati[0] != Bust {
		t.Errorf("Expected AI player to stand or bust, but got %s", aiPlayer.Stati[0])
	}
}

func TestDetermineWinner(t *testing.T) {
	testCases := []struct {
		player      Player
		dealer      Player
		expectedPlayerStatus PlayerStatus
		expectedPlayerChips int
	}{
		{
			player:      Player{Name: "Player 1", Hands: []Hand{{}}, Scores: []int{20}, Stati: []PlayerStatus{Stand}, Bets: []int{100}, IsHuman: true, Chips: 900},
			dealer:      Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			expectedPlayerStatus: PlayerWins,
			expectedPlayerChips: 1100,
		},
		{
			player:      Player{Name: "Player 1", Hands: []Hand{{}}, Scores: []int{18}, Stati: []PlayerStatus{Stand}, Bets: []int{100}, IsHuman: true, Chips: 900},
			dealer:      Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			expectedPlayerStatus: DealerWins,
			expectedPlayerChips: 900,
		},
		{
			player:      Player{Name: "Player 1", Hands: []Hand{{}}, Scores: []int{19}, Stati: []PlayerStatus{Stand}, Bets: []int{100}, IsHuman: true, Chips: 900},
			dealer:      Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			expectedPlayerStatus: Push,
			expectedPlayerChips: 1000,
		},
		{
			player:      Player{Name: "Player 1", Hands: []Hand{{}}, Scores: []int{22}, Stati: []PlayerStatus{Bust}, Bets: []int{100}, IsHuman: true, Chips: 900},
			dealer:      Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			expectedPlayerStatus: DealerWins,
			expectedPlayerChips: 900,
		},
		{
			player:      Player{Name: "Player 1", Hands: []Hand{{}}, Scores: []int{21}, Stati: []PlayerStatus{BlackjackWin}, Bets: []int{100}, IsHuman: true, Chips: 900},
			dealer:      Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			expectedPlayerStatus: BlackjackWin,
			expectedPlayerChips: 1150,
		},
	}

	for _, tc := range testCases {
		game := NewGame()
		game.Players[0] = tc.player
		game.Dealer = tc.dealer
		game.determineWinner()
		if game.Players[0].Stati[0] != tc.expectedPlayerStatus {
			t.Errorf("Expected player status to be %s, but got %s", tc.expectedPlayerStatus, game.Players[0].Stati[0])
		}
		if game.Players[0].Chips != tc.expectedPlayerChips {
			t.Errorf("Expected player chips to be %d, but got %d", tc.expectedPlayerChips, game.Players[0].Chips)
		}
	}
}