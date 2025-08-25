const messageBar = document.getElementById("message-bar");
const dealerCardsEl = document.getElementById("dealer-cards");
const dealerScoreEl = document.getElementById("dealer-score");
const playerCardsEl = document.getElementById("player-cards");
const playerScoreEl = document.getElementById("player-score");
const playerChipsEl = document.getElementById("player-chips");

const newGameBtn = document.getElementById("new-game-btn");
const hitBtn = document.getElementById("hit-btn");
const standBtn = document.getElementById("stand-btn");
const betBtn = document.getElementById("bet-btn");
const betAmountInput = document.getElementById("bet-amount");
const doubleDownBtn = document.getElementById("double-down-btn");

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


    // Render the player's hand
    playerCardsEl.innerHTML = "";
    if (data.Player && data.Player.Hands && data.Player.Hands[0]) {
        for (const card of data.Player.Hands[0]) {
            const cardEl = document.createElement("div");
            cardEl.classList.add("card");
            cardEl.textContent = getCardName(card);
            playerCardsEl.appendChild(cardEl);
        }
        playerScoreEl.textContent = data.Player.Scores[0];
    } else {
        playerScoreEl.textContent = "";
    }

    playerChipsEl.textContent = data.PlayerChips;

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
        const playerStatus = data.Player.Stati[0];
        switch (playerStatus) {
            case "player_wins":
                messageBar.textContent = "You win!";
                break;
            case "dealer_wins":
                messageBar.textContent = "Dealer wins!";
                break;
            case "push":
                messageBar.textContent = "It's a push!";
                break;
            case "bust":
                messageBar.textContent = "Bust!";
                break;
            default:
                messageBar.textContent = "Game over! Place your bet for the next round.";
        }
    }

    // Show/hide action buttons
    hitBtn.style.display = data.AvailableActions.includes("hit") ? "inline-block" : "none";
    standBtn.style.display = data.AvailableActions.includes("stand") ? "inline-block" : "none";
    doubleDownBtn.style.display = data.AvailableActions.includes("doubledown") ? "inline-block" : "none";
    bettingControls.style.display = data.AvailableActions.includes("bet") ? "block" : "none";
}

newGameBtn.addEventListener("click", function() {
    fetch("/api/new_game")
        .then(response => response.json())
        .then(data => {
            renderGame(data);
        });
});

betBtn.addEventListener("click", function() {
    const amount = betAmountInput.value;
    fetch(`/api/bet?amount=${amount}`)
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

doubleDownBtn.addEventListener("click", function() {
    fetch("/api/doubledown")
        .then(response => response.json())
        .then(data => {
            renderGame(data);
        });
});

// Initial game load
fetch("/api/new_game")
    .then(response => response.json())
    .then(data => {
        renderGame(data);
    });