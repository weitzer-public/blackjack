package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var game Game

func main() {
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/api/new_game", newGameHandler)
	http.HandleFunc("/api/hit", hitHandler)
	http.HandleFunc("/api/stand", standHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	game = NewGame()
	visibleGame := VisibleGame{
		Player:      game.Player,
		Dealer:      []Card{game.Dealer[0]}, // Only show one card
		PlayerScore: game.PlayerScore,
		DealerScore: HandScore(Hand{game.Dealer[0]}),
		State:       game.State,
	}
	json.NewEncoder(w).Encode(visibleGame)
}

func hitHandler(w http.ResponseWriter, r *http.Request) {
	if game.State != "playing" {
		json.NewEncoder(w).Encode(game)
		return
	}
	game.Player = append(game.Player, game.Deck[0])
	game.Deck = game.Deck[1:]
	game.PlayerScore = HandScore(game.Player)
	if game.PlayerScore > 21 {
		game.State = "player_busts"
	}

	visibleGame := VisibleGame{
		Player:      game.Player,
		Dealer:      []Card{game.Dealer[0]},
		PlayerScore: game.PlayerScore,
		DealerScore: HandScore(Hand{game.Dealer[0]}),
		State:       game.State,
	}
	json.NewEncoder(w).Encode(visibleGame)
}

func standHandler(w http.ResponseWriter, r *http.Request) {
	if game.State != "playing" {
		json.NewEncoder(w).Encode(game)
		return
	}
	// Dealer's turn
	for HandScore(game.Dealer) < 17 {
		game.Dealer = append(game.Dealer, game.Deck[0])
		game.Deck = game.Deck[1:]
	}
	game.DealerScore = HandScore(game.Dealer)

	if game.DealerScore > 21 || game.PlayerScore > game.DealerScore {
		game.State = "player_wins"
	} else if game.PlayerScore < game.DealerScore {
		game.State = "dealer_wins"
	} else {
		game.State = "tie"
	}

	json.NewEncoder(w).Encode(game)
}