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
    switch (data.GameState) {
        case "betting":
            messageBar.textContent = "Place your bet!";
            bettingControls.style.display = "block";
            gameControls.style.display = "none";
            break;
        case "playing":
            messageBar.textContent = "Your turn!";
            bettingControls.style.display = "none";
            gameControls.style.display = "block";
            break;
        case "game_over":
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

// Initial game load
fetch("/api/new_game")
    .then(response => response.json())
    .then(data => {
        renderGame(data);
    });