// You would need to load game.js in your HTML before this script
// or use a module bundler. For simplicity, this example assumes game.js is available.

const messageEl = document.getElementById("message-el");
const sumEl = document.getElementById("sum-el");
const cardsEl = document.getElementById("cards-el");
const playerEl = document.getElementById("player-el");

// Initial player display
let gameState = game.getGameState();
playerEl.textContent = `${gameState.player.name}: $${gameState.player.chips}`;

function renderGame(state) {
    cardsEl.textContent = "Cards: ";
    for (let i = 0; i < state.cards.length; i++) {
        cardsEl.textContent += state.cards[i] + " ";
    }
    
    sumEl.textContent = "Sum: " + state.sum;
    messageEl.textContent = state.message;
}

// Example of how you'd hook up the UI
document.querySelector("#start-game-btn").addEventListener("click", () => {
    const state = game.startGame();
    renderGame(state);
});

document.querySelector("#new-card-btn").addEventListener("click", () => {
    const state = game.newCard();
    if (state) { // newCard only returns a state if a card was drawn
        renderGame(state);
    }
});