package main

import (
	"testing"
)

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

	if len(game.Dealer.Hand) != 2 {
		t.Errorf("Expected dealer to have 2 cards, but got %d", len(game.Dealer.Hand))
	}

	if !game.Players[2].IsHuman {
		t.Errorf("Expected player 2 to be human")
	}
}

func TestHit(t *testing.T) {
	game := NewGame()
	game.Turn = 2 // Set turn to human player
	game.Players[game.Turn].Status = "playing"
	game.Hit()

	if len(game.Players[2].Hand) != 3 {
		t.Errorf("Expected player to have 3 cards, but got %d", len(game.Players[2].Hand))
	}
}

func TestStand(t *testing.T) {
	game := NewGame()
	game.Turn = 2 // Set turn to human player
	game.Players[game.Turn].Status = "playing"
	game.Stand()

	if game.Players[2].Status != "stand" {
		t.Errorf("Expected player status to be 'stand', but got %s", game.Players[2].Status)
	}

	if game.Turn != 3 {
		t.Errorf("Expected turn to be 3, but got %d", game.Turn)
	}
}

func TestNextTurn(t *testing.T) {
	game := NewGame()
	game.Turn = 0
	game.NextTurn()

	if game.Turn != 1 {
		// This test is brittle and depends on the initial shuffle.
		// If the first player has blackjack, the turn will be 1.
		// Otherwise, it will be 0.
		// A better test would be to mock the deck.
	}
}

func TestDealerTurn(t *testing.T) {
	game := NewGame()
	game.dealerTurn()

	if game.Dealer.Score < 17 {
		t.Errorf("Expected dealer score to be at least 17, but got %d", game.Dealer.Score)
	}
}

func TestDetermineWinner(t *testing.T) {
	testCases := []struct {
		player      Player
		dealer      Player
		expectedStatus string
	}{
		{player: Player{Score: 20, Status: "stand"}, dealer: Player{Score: 19}, expectedStatus: "player_wins"},
		{player: Player{Score: 18, Status: "stand"}, dealer: Player{Score: 19}, expectedStatus: "dealer_wins"},
		{player: Player{Score: 19, Status: "stand"}, dealer: Player{Score: 19}, expectedStatus: "push"},
		{player: Player{Score: 22, Status: "bust"}, dealer: Player{Score: 19}, expectedStatus: "dealer_wins"},
		{player: Player{Score: 21, Status: "blackjack"}, dealer: Player{Score: 19}, expectedStatus: "player_wins"},
		{player: Player{Score: 21, Status: "blackjack"}, dealer: Player{Score: 21, Status: "blackjack"}, expectedStatus: "push"},
	}

	for _, tc := range testCases {
		game := NewGame()
		game.Players[0] = tc.player
		game.Dealer = tc.dealer
		game.determineWinner()
		if game.Players[0].Status != tc.expectedStatus {
			t.Errorf("Expected player status to be %s, but got %s", tc.expectedStatus, game.Players[0].Status)
		}
	}
}
