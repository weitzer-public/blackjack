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

	if len(game.Player) != 2 {
		t.Errorf("Expected player to have 2 cards, but got %d", len(game.Player))
	}

	if len(game.Dealer) != 1 {
		t.Errorf("Expected dealer to have 1 card, but got %d", len(game.Dealer))
	}

	if game.State != "playing" {
		t.Errorf("Expected game state to be 'playing', but got %s", game.State)
	}
}

func TestHit(t *testing.T) {
	game := NewGame()
	game.Hit()

	if len(game.Player) != 3 {
		t.Errorf("Expected player to have 3 cards, but got %d", len(game.Player))
	}
}

func TestStand(t *testing.T) {
	game := NewGame()
	game.Stand()

	if game.State == "playing" {
		t.Errorf("Expected game state to not be 'playing', but got %s", game.State)
	}
}