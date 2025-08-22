package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

var (
	game  Game
	mutex = &sync.Mutex{}
)

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
	mutex.Lock()
	defer mutex.Unlock()
	game = NewGame()
	game.NextTurn()
	json.NewEncoder(w).Encode(game.Visible())
}

func hitHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	if game.Players[game.Turn].IsHuman {
		game.Hit()
		if game.Players[game.Turn].Status == Bust {
			game.NextTurn()
		}
	}
	json.NewEncoder(w).Encode(game.Visible())
}

func standHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	if game.Players[game.Turn].IsHuman {
		game.Stand()
		game.NextTurn()
	}
	json.NewEncoder(w).Encode(game.Visible())
}
