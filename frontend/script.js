const messageEl = document.getElementById("message-el");
const dealerCardsEl = document.getElementById("dealer-cards");
const dealerScoreEl = document.getElementById("dealer-score");
const playersContainerEl = document.getElementById("players-container");
const newGameBtn = document.getElementById("new-game-btn");
const hitBtn = document.getElementById("hit-btn");
const standBtn = document.getElementById("stand-btn");

function getCardName(card) {
    const suits = ["♠", "♥", "♦", "♣"];
    const values = ["A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"];
    return `${values[card.Value - 1]}${suits[card.Suit]}`;
}

function renderGame(data) {
    // Render the dealer's hand
    dealerCardsEl.innerHTML = "";
    for (const card of data.Dealer.Hand) {
        const cardEl = document.createElement("div");
        cardEl.classList.add("card");
        cardEl.textContent = getCardName(card);
        dealerCardsEl.appendChild(cardEl);
    }
    dealerScoreEl.textContent = data.Dealer.Score;

    // Render the players' hands
    playersContainerEl.innerHTML = "";
    for (const player of data.Players) {
        const playerEl = document.createElement("div");
        playerEl.classList.add("player");
        if (player.IsHuman) {
            playerEl.classList.add("human");
        }

        const playerNameEl = document.createElement("h3");
        playerNameEl.textContent = player.IsHuman ? "You" : "Computer";
        playerEl.appendChild(playerNameEl);

        const cardsEl = document.createElement("div");
        cardsEl.classList.add("cards");
        for (const card of player.Hand) {
            const cardEl = document.createElement("div");
            cardEl.classList.add("card");
            cardEl.textContent = getCardName(card);
            cardsEl.appendChild(cardEl);
        }
        playerEl.appendChild(cardsEl);

        const scoreEl = document.createElement("p");
        scoreEl.textContent = "Score: " + player.Score;
        playerEl.appendChild(scoreEl);

        const statusEl = document.createElement("p");
        statusEl.textContent = "Status: " + player.Status;
        playerEl.appendChild(statusEl);

        playersContainerEl.appendChild(playerEl);
    }

    // Display the game message
    switch (data.State) {
        case "playing":
            messageEl.textContent = "Your turn!";
            break;
        case "game_over":
            const humanPlayer = data.Players.find(p => p.IsHuman);
            switch (humanPlayer.Status) {
                case "player_wins":
                    messageEl.textContent = "You win!";
                    break;
                case "dealer_wins":
                    messageEl.textContent = "Dealer wins!";
                    break;
                case "push":
                    messageEl.textContent = "It's a push!";
                    break;
                case "bust":
                    messageEl.textContent = "Bust!";
                    break;
                default:
                    messageEl.textContent = "Game over!";
            }
            break;
    }
}

newGameBtn.addEventListener("click", function() {
    fetch("/api/new_game")
        .then(response => response.json())
        .then(data => {
            renderGame(data);
        });
});

hitBtn.addEventListener("click", function() {
    fetch("/api/hit")
        .then(response => response.json())
        .then(data => {
            renderGame(data);
        });
});

standBtn.addEventListener("click", function() {
    fetch("/api/stand")
        .then(response => response.json())
        .then(data => {
            renderGame(data);
        });
});
