package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func betHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	amount, err := strconv.Atoi(r.URL.Query().Get("amount"))
	if err != nil {
		http.Error(w, "Invalid bet amount", http.StatusBadRequest)
		return
	}

	game.PlaceBet(amount)
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