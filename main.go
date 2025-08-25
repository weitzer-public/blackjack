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
	http.HandleFunc("/api/bet", betHandler)
	http.HandleFunc("/api/hit", hitHandler)
	http.HandleFunc("/api/stand", standHandler)
	http.HandleFunc("/api/doubledown", doubleDownHandler)
	http.HandleFunc("/api/split", splitHandler)
	http.HandleFunc("/api/game", gameHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	game = NewGame()
	json.NewEncoder(w).Encode(game.Visible())
}

type BetRequest struct {
	Amount int `json:"amount"`
}

func betHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var req BetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	game.PlaceBet(req.Amount)
	json.NewEncoder(w).Encode(game.Visible())
}

func hitHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	game.Hit()
	json.NewEncoder(w).Encode(game.Visible())
}

func standHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	game.Stand()
	json.NewEncoder(w).Encode(game.Visible())
}

func doubleDownHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	game.DoubleDown()
	json.NewEncoder(w).Encode(game.Visible())
}

func splitHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	game.Split()
	json.NewEncoder(w).Encode(game.Visible())
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	json.NewEncoder(w).Encode(game.Visible())
}
