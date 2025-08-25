const messageBar = document.getElementById("message-bar");
const dealerCardsEl = document.getElementById("dealer-cards");
const dealerScoreEl = document.getElementById("dealer-score");
const playerCardsEl = document.getElementById("player-cards");
const playerScoreEl = document.getElementById("player-score");
const playersArea = document.getElementById("players-area");
const playerChipsEl = document.getElementById("player-chips");

const newGameBtn = document.getElementById("new-game-btn");
const hitBtn = document.getElementById("hit-btn");
const standBtn = document.getElementById("stand-btn");
const betBtn = document.getElementById("bet-btn");
const betAmountInput = document.getElementById("bet-amount");
const doubleDownBtn = document.getElementById("double-down-btn");
const splitBtn = document.getElementById("split-btn");

const bettingControls = document.getElementById("betting-controls");
const gameControls = document.getElementById("game-controls");

function getCardName(card) {
    const suits = ["♠", "♥", "♦", "♣"];
    const values = ["A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"];
    return `${values[card.Value - 1]}${suits[card.Suit]}`;
}

function renderGame(data) {
    // Render the dealer's hand
    dealerCardsEl.innerHTML = "";
    if (data.Dealer && data.Dealer.Hands && data.Dealer.Hands[0]) {
        for (const card of data.Dealer.Hands[0]) {
            const cardEl = document.createElement("div");
            cardEl.classList.add("card");
            cardEl.textContent = getCardName(card);
            dealerCardsEl.appendChild(cardEl);
        }
        dealerScoreEl.textContent = data.Dealer.Scores[0];
    } else {
        dealerScoreEl.textContent = "";
    }

    // Render AI players
    playersArea.innerHTML = "";
    if(data.Players) {
        for (const player of data.Players) {
            if (player.IsHuman) {
                renderPlayer(player, true);
            } else {
                renderPlayer(player, false);
            }
        }
    }


    // Update UI based on game state
    if (data.GameState === "betting") {
        messageBar.textContent = "Place your bet!";
        bettingControls.style.display = "block";
        gameControls.style.display = "none";
    } else if (data.GameState === "playing") {
        messageBar.textContent = "Your turn!";
        bettingControls.style.display = "none";
        gameControls.style.display = "block";
    } else if (data.GameState === "game_over") {
        bettingControls.style.display = "block"; // Allow betting for next game
        gameControls.style.display = "none";
        messageBar.textContent = "Game over! Place your bet for the next round.";
    }

    // Show/hide action buttons
    hitBtn.style.display = data.AvailableActions.includes("hit") ? "inline-block" : "none";
    standBtn.style.display = data.AvailableActions.includes("stand") ? "inline-block" : "none";
    doubleDownBtn.style.display = data.AvailableActions.includes("doubledown") ? "inline-block" : "none";
    splitBtn.style.display = data.AvailableActions.includes("split") ? "inline-block" : "none";
    bettingControls.style.display = data.AvailableActions.includes("bet") ? "block" : "none";
}

function renderPlayer(player, isHuman) {
    let playerEl;
    if(isHuman) {
        playerEl = document.getElementById("player-area");
        playerChipsEl.textContent = player.Chips;
    } else {
        playerEl = document.createElement("div");
        playerEl.classList.add("player");
        playersArea.appendChild(playerEl);
    }

    const nameEl = document.createElement("h2");
    nameEl.textContent = player.Name;
    playerEl.appendChild(nameEl);

    const cardsEl = document.createElement("div");
    cardsEl.classList.add("cards");
    if (player.Hands) {
        for (let i = 0; i < player.Hands.length; i++) {
            const hand = player.Hands[i];
            const handEl = document.createElement("div");
            handEl.classList.add("hand");

            const handCardsEl = document.createElement("div");
            handCardsEl.classList.add("cards");
            for (const card of hand) {
                const cardEl = document.createElement("div");
                cardEl.classList.add("card");
                cardEl.textContent = getCardName(card);
                handCardsEl.appendChild(cardEl);
            }
            handEl.appendChild(handCardsEl);

            const scoreEl = document.createElement("p");
            scoreEl.textContent = "Score: " + player.Scores[i];
            handEl.appendChild(scoreEl);

            const statusEl = document.createElement("p");
            statusEl.textContent = "Status: " + player.Stati[i];
            handEl.appendChild(statusEl);

            cardsEl.appendChild(handEl);
        }
    }
    playerEl.appendChild(cardsEl);
}

function performAction(url, method = 'GET', data = null) {
    const options = {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        }
    };

    if (data && method === 'POST') {
        options.body = JSON.stringify(data);
    }

    fetch(url, options)
    .then(response => {
        if (!response.ok) {
            return response.json().then(err => {
                throw new Error(err.error || 'Server error');
            });
        }
        return response.json();
    })
    .then(data => {
        renderGame(data);
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
        messageBar.textContent = "Error: " + error.message;
    });
}

newGameBtn.addEventListener("click", function() {
    performAction("/api/new_game", 'POST');
});

betBtn.addEventListener("click", function() {
    const amount = parseInt(betAmountInput.value, 10);
    performAction('/api/bet', 'POST', { amount });
});

hitBtn.addEventListener("click", function() {
    performAction("/api/hit", 'POST');
});

standBtn.addEventListener("click", function() {
    performAction("/api/stand", 'POST');
});

doubleDownBtn.addEventListener("click", function() {
    performAction("/api/doubledown", 'POST');
});

splitBtn.addEventListener("click", function() {
    performAction("/api/split", 'POST');
});

// Initial game load
performAction("/api/game", 'GET');