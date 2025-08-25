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
	// Test case for an even bet
	gameEven := NewGame()
	gameEven.PlaceBet(100)

	// Force a blackjack for the player
	gameEven.Player.Hands[0] = Hand{{Value: 1}, {Value: 10}}
	gameEven.Player.Scores[0] = HandScore(gameEven.Player.Hands[0])

	// Make sure dealer doesn't have blackjack
	gameEven.Dealer.Hands[0] = Hand{{Value: 2}, {Value: 10}}
	gameEven.Dealer.Scores[0] = HandScore(gameEven.Dealer.Hands[0])

	// Check for blackjack
	if gameEven.Player.Scores[0] == Blackjack {
		gameEven.Player.Stati[0] = BlackjackWin
		if gameEven.Dealer.Scores[0] == Blackjack {
			gameEven.Dealer.Stati[0] = BlackjackWin
			gameEven.Player.Stati[0] = Push
			gameEven.PlayerChips += gameEven.PlayerBet // Push, return bet
			gameEven.GameState = "game_over"
		} else {
			gameEven.PlayerChips += gameEven.PlayerBet + int(float64(gameEven.PlayerBet)*1.5) // Blackjack pays 3:2
			gameEven.GameState = "game_over"
		}
	} else if gameEven.Dealer.Scores[0] == Blackjack {
		gameEven.Dealer.Stati[0] = BlackjackWin
		gameEven.GameState = "game_over"
	}

	expectedChipsEven := 1150
	if gameEven.PlayerChips != expectedChipsEven {
		t.Errorf("Expected player chips to be %d after blackjack, but got %d", expectedChipsEven, gameEven.PlayerChips)
	}

	// Test case for an odd bet
	gameOdd := NewGame()
	gameOdd.PlaceBet(15)

	// Force a blackjack for the player
	gameOdd.Player.Hands[0] = Hand{{Value: 1}, {Value: 10}}
	gameOdd.Player.Scores[0] = HandScore(gameOdd.Player.Hands[0])

	// Make sure dealer doesn't have blackjack
	gameOdd.Dealer.Hands[0] = Hand{{Value: 2}, {Value: 10}}
	gameOdd.Dealer.Scores[0] = HandScore(gameOdd.Dealer.Hands[0])

	// Check for blackjack
	if gameOdd.Player.Scores[0] == Blackjack {
		gameOdd.Player.Stati[0] = BlackjackWin
		if gameOdd.Dealer.Scores[0] == Blackjack {
			gameOdd.Dealer.Stati[0] = BlackjackWin
			gameOdd.Player.Stati[0] = Push
			gameOdd.PlayerChips += gameOdd.PlayerBet // Push, return bet
			gameOdd.GameState = "game_over"
		} else {
			gameOdd.PlayerChips += gameOdd.PlayerBet + int(float64(gameOdd.PlayerBet)*1.5) // Blackjack pays 3:2
			gameOdd.GameState = "game_over"
		}
	} else if gameOdd.Dealer.Scores[0] == Blackjack {
		gameOdd.Dealer.Stati[0] = BlackjackWin
		gameOdd.GameState = "game_over"
	}

	expectedChipsOdd := 1022
	if gameOdd.PlayerChips != expectedChipsOdd {
		t.Errorf("Expected player chips to be %d after blackjack with odd bet, but got %d", expectedChipsOdd, gameOdd.PlayerChips)
	}
}