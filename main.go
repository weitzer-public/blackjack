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
	json.NewEncoder(w).Encode(game)
}

func hitHandler(w http.ResponseWriter, r *http.Request) {
	game.Player = append(game.Player, game.Deck[0])
	game.Deck = game.Deck[1:]
	json.NewEncoder(w).Encode(game)
}

func standHandler(w http.ResponseWriter, r *http.Request) {
	// Dealer's turn
	for HandScore(game.Dealer) < 17 {
		game.Dealer = append(game.Dealer, game.Deck[0])
		game.Deck = game.Deck[1:]
	}
	json.NewEncoder(w).Encode(game)
}

