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

	if len(game.Players) != NumAIPlayers+1 {
		t.Errorf("Expected %d players, but got %d", NumAIPlayers+1, len(game.Players))
	}

	if game.GameState != "betting" {
		t.Errorf("Expected game state to be 'betting', but got %s", game.GameState)
	}
}

func TestPlaceBet(t *testing.T) {
	game := NewGame()
	game.PlaceBet(100)

	if game.Players[0].Chips != 900 {
		t.Errorf("Expected player to have 900 chips, but got %d", game.Players[0].Chips)
	}

	if game.Players[0].Bets[0] != 100 {
		t.Errorf("Expected player bet to be 100, but got %d", game.Players[0].Bets[0])
	}

	for i := 1; i <= NumAIPlayers; i++ {
		if game.Players[i].Bets[0] != AIPlayerBet {
			t.Errorf("Expected AI player %d bet to be %d, but got %d", i, AIPlayerBet, game.Players[i].Bets[0])
		}
	}

	if game.GameState != "playing" {
		t.Errorf("Expected game state to be 'playing', but got %s", game.GameState)
	}

	if len(game.Players[0].Hands[0]) != 2 {
		t.Errorf("Expected player to have 2 cards, but got %d", len(game.Players[0].Hands[0]))
	}

	if len(game.Dealer.Hands[0]) != 2 {
		t.Errorf("Expected dealer to have 2 cards, but got %d", len(game.Dealer.Hands[0]))
	}
}

func TestHit(t *testing.T) {
	game := NewGame()
	game.PlaceBet(100)
	game.Hit()

	if len(game.Players[0].Hands[0]) != 3 {
		t.Errorf("Expected player to have 3 cards, but got %d", len(game.Players[0].Hands[0]))
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

func TestDealerTurn(t *testing.T) {
	game := NewGame()
	game.PlaceBet(100)
	game.dealerTurn()

	if game.Dealer.Scores[0] < 17 {
		t.Errorf("Expected dealer score to be at least 17, but got %d", game.Dealer.Scores[0])
	}
}

func TestDetermineWinner(t *testing.T) {
	testCases := []struct {
		name                string
		player              Player
		dealer              Player
		expectedPlayerStatus PlayerStatus
		expectedPlayerChips int
	}{
		{
			name:                "Player wins",
			player:              Player{Hands: []Hand{{}}, Scores: []int{20}, Stati: []PlayerStatus{Stand}, Bets: []int{100}, Chips: 900},
			dealer:              Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			expectedPlayerStatus: PlayerWins,
			expectedPlayerChips: 1100,
		},
		{
			name:                "Dealer wins",
			player:              Player{Hands: []Hand{{}}, Scores: []int{18}, Stati: []PlayerStatus{Stand}, Bets: []int{100}, Chips: 900},
			dealer:              Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			expectedPlayerStatus: DealerWins,
			expectedPlayerChips: 900,
		},
		{
			name:                "Push",
			player:              Player{Hands: []Hand{{}}, Scores: []int{19}, Stati: []PlayerStatus{Stand}, Bets: []int{100}, Chips: 900},
			dealer:              Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			expectedPlayerStatus: Push,
			expectedPlayerChips: 1000,
		},
		{
			name:                "Player busts",
			player:              Player{Hands: []Hand{{}}, Scores: []int{22}, Stati: []PlayerStatus{Bust}, Bets: []int{100}, Chips: 900},
			dealer:              Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			expectedPlayerStatus: DealerWins,
			expectedPlayerChips: 900,
		},
		{
			name:                "Dealer busts",
			player:              Player{Hands: []Hand{{}}, Scores: []int{20}, Stati: []PlayerStatus{Stand}, Bets: []int{100}, Chips: 900},
			dealer:              Player{Scores: []int{22}, Stati: []PlayerStatus{Bust}},
			expectedPlayerStatus: PlayerWins,
			expectedPlayerChips: 1100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := NewGame()
			game.Players = []Player{tc.player}
			game.Dealer = tc.dealer
			game.determineWinner()
			if game.Players[0].Stati[0] != tc.expectedPlayerStatus {
				t.Errorf("Expected player status to be %s, but got %s", tc.expectedPlayerStatus, game.Players[0].Stati[0])
			}
			if game.Players[0].Chips != tc.expectedPlayerChips {
				t.Errorf("Expected player chips to be %d, but got %d", tc.expectedPlayerChips, game.Players[0].Chips)
			}
		})
	}
}

func TestBlackjackPayout(t *testing.T) {
	game := NewGame()
	game.PlaceBet(100)

	// Force a blackjack for the player
	game.Players[0].Hands[0] = Hand{{Value: 1}, {Value: 10}}
	game.Players[0].Scores[0] = HandScore(game.Players[0].Hands[0])
	game.Players[0].Stati[0] = BlackjackWin

	// Make sure dealer doesn't have blackjack
	game.Dealer.Hands[0] = Hand{{Value: 2}, {Value: 10}}
	game.Dealer.Scores[0] = HandScore(game.Dealer.Hands[0])

	game.determineWinner()

	// Player wins with blackjack, payout should be 3:2
	// Initial chips: 1000, Bet: 100. Chips after bet: 900.
	// Payout: 100 (original bet) + 100 * 3 / 2 = 150. Total payout: 250
	// Expected chips: 900 + 250 = 1150
	expectedChips := 1150
	if game.Players[0].Chips != expectedChips {
		t.Errorf("Expected player chips to be %d after blackjack, but got %d", expectedChips, game.Players[0].Chips)
	}
}