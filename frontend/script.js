const messageEl = document.getElementById("message-el");
const dealerCardsEl = document.getElementById("dealer-cards");
const dealerScoreEl = document.getElementById("dealer-score");
const playersContainerEl = document.getElementById("players-container");
const newGameBtn = document.getElementById("new-game-btn");
const hitBtn = document.getElementById("hit-btn");
const standBtn = document.getElementById("stand-btn");

function renderGame(data) {
    // Render the dealer's hand
    dealerCardsEl.innerHTML = "";
    for (const card of data.Dealer.Hand) {
        const cardEl = document.createElement("div");
        cardEl.classList.add("card");
        cardEl.textContent = card.Value;
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
            cardEl.textContent = card.Value;
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
            messageEl.textContent = "Game over!";
            break;
        default:
            messageEl.textContent = data.State;
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
