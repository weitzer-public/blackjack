package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var game Game

func main() {
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/api/new_game", newGameHandler)
	http.HandleFunc("/api/hit", hitHandler)
	http.HandleFunc("/api/stand", standHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	game = NewGame()
	visibleGame := VisibleGame{
		Player:      game.Player,
		Dealer:      game.Dealer,
		PlayerScore: game.PlayerScore,
		DealerScore: game.DealerScore,
		State:       game.State,
	}
	json.NewEncoder(w).Encode(visibleGame)
}

func hitHandler(w http.ResponseWriter, r *http.Request) {
	game.Hit()
	visibleGame := VisibleGame{
		Player:      game.Player,
		Dealer:      game.Dealer,
		PlayerScore: game.PlayerScore,
		DealerScore: game.DealerScore,
		State:       game.State,
	}
	json.NewEncoder(w).Encode(visibleGame)
}

func standHandler(w http.ResponseWriter, r *http.Request) {
	game.Stand()
	json.NewEncoder(w).Encode(game)
}