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

	if game.PlayerChips != 1000 {
		t.Errorf("Expected player to have 1000 chips, but got %d", game.PlayerChips)
	}

	if game.GameState != "betting" {
		t.Errorf("Expected game state to be 'betting', but got %s", game.GameState)
	}
}

func TestPlaceBet(t *testing.T) {
	game := NewGame()
	game.PlaceBet(100)

	if game.PlayerChips != 900 {
		t.Errorf("Expected player to have 900 chips, but got %d", game.PlayerChips)
	}

	if game.PlayerBet != 100 {
		t.Errorf("Expected player bet to be 100, but got %d", game.PlayerBet)
	}

	if game.GameState != "playing" {
		t.Errorf("Expected game state to be 'playing', but got %s", game.GameState)
	}

	if len(game.Player.Hands[0]) != 2 {
		t.Errorf("Expected player to have 2 cards, but got %d", len(game.Player.Hands[0]))
	}

	if len(game.Dealer.Hands[0]) != 2 {
		t.Errorf("Expected dealer to have 2 cards, but got %d", len(game.Dealer.Hands[0]))
	}
}

func TestHit(t *testing.T) {
	game := NewGame()
	game.PlaceBet(100)
	game.Hit()

	if len(game.Player.Hands[0]) != 3 {
		t.Errorf("Expected player to have 3 cards, but got %d", len(game.Player.Hands[0]))
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
		player      Player
		dealer      Player
		playerChips int
		expectedPlayerStatus PlayerStatus
		expectedPlayerChips int
	}{
		{
			player:      Player{Hands: []Hand{{}}, Scores: []int{20}, Stati: []PlayerStatus{Stand}, Bets: []int{100}},
			dealer:      Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			playerChips: 900,
			expectedPlayerStatus: PlayerWins,
			expectedPlayerChips: 1100,
		},
		{
			player:      Player{Hands: []Hand{{}}, Scores: []int{18}, Stati: []PlayerStatus{Stand}, Bets: []int{100}},
			dealer:      Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			playerChips: 900,
			expectedPlayerStatus: DealerWins,
			expectedPlayerChips: 900,
		},
		{
			player:      Player{Hands: []Hand{{}}, Scores: []int{19}, Stati: []PlayerStatus{Stand}, Bets: []int{100}},
			dealer:      Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			playerChips: 900,
			expectedPlayerStatus: Push,
			expectedPlayerChips: 1000,
		},
		{
			player:      Player{Hands: []Hand{{}}, Scores: []int{22}, Stati: []PlayerStatus{Bust}, Bets: []int{100}},
			dealer:      Player{Scores: []int{19}, Stati: []PlayerStatus{Stand}},
			playerChips: 900,
			expectedPlayerStatus: DealerWins,
			expectedPlayerChips: 900,
		},
	}

	for _, tc := range testCases {
		game := NewGame()
		game.Player = tc.player
		game.Dealer = tc.dealer
		game.PlayerChips = tc.playerChips
		game.PlayerBet = 100
		game.determineWinner()
		if game.Player.Stati[0] != tc.expectedPlayerStatus {
			t.Errorf("Expected player status to be %s, but got %s", tc.expectedPlayerStatus, game.Player.Stati[0])
		}
		if game.PlayerChips != tc.expectedPlayerChips {
			t.Errorf("Expected player chips to be %d, but got %d", tc.expectedPlayerChips, game.PlayerChips)
		}
	}
}

func TestBlackjackPayout(t *testing.T) {
	game := NewGame()
	game.PlaceBet(100)

	// Force a blackjack for the player
	game.Player.Hands[0] = Hand{{Value: 1}, {Value: 10}}
	game.Player.Scores[0] = HandScore(game.Player.Hands[0])

	// Make sure dealer doesn't have blackjack
	game.Dealer.Hands[0] = Hand{{Value: 2}, {Value: 10}}
	game.Dealer.Scores[0] = HandScore(game.Dealer.Hands[0])

	// Check for blackjack
	if game.Player.Scores[0] == Blackjack {
		game.Player.Stati[0] = BlackjackWin
		if game.Dealer.Scores[0] == Blackjack {
			game.Dealer.Stati[0] = BlackjackWin
			game.Player.Stati[0] = Push
			game.PlayerChips += game.PlayerBet // Push, return bet
			game.GameState = "game_over"
		} else {
			game.PlayerChips += game.PlayerBet + (game.PlayerBet*3)/2 // Blackjack pays 3:2
			game.GameState = "game_over"
		}
	} else if game.Dealer.Scores[0] == Blackjack {
		game.Dealer.Stati[0] = BlackjackWin
		game.GameState = "game_over"
	}

	// Player wins with blackjack, payout should be 3:2
	// Initial chips: 1000, Bet: 100. Chips after bet: 900.
	// Payout: 100 (original bet) + 100 * 3 / 2 = 150. Total payout: 250
	// Expected chips: 900 + 250 = 1150
	expectedChips := 1150
	if game.PlayerChips != expectedChips {
		t.Errorf("Expected player chips to be %d after blackjack, but got %d", expectedChips, game.PlayerChips)
	}
}
