const messageEl = document.getElementById("message-el")
const sumEl = document.getElementById("sum-el")
const cardsEl = document.getElementById("cards-el")
const playerEl = document.getElementById("player-el")
const startGameBtn = document.getElementById("start-game-btn")
const newCardBtn = document.getElementById("new-card-btn")

let player = {
    name: "Per",
    chips: 200
}

playerEl.textContent = `${player.name}: ${player.chips}`

function renderGame(data) {
    cardsEl.textContent = "Cards: "
    for (let i = 0; i < data.Player.length; i++) {
        cardsEl.textContent += data.Player[i].Value + " "
    }

    let playerScore = 0;
    for (let i = 0; i < data.Player.length; i++) {
        if (data.Player[i].Value > 10) {
            playerScore += 10;
        } else {
            playerScore += data.Player[i].Value;
        }
    }

    sumEl.textContent = "Sum: " + playerScore

    if (playerScore <= 20) {
        messageEl.textContent = "Do you want to draw a new card?"
    } else if (playerScore === 21) {
        messageEl.textContent = "You've got Blackjack!"
    } else {
        messageEl.textContent = "You're out of the game!"
    }
}

startGameBtn.addEventListener("click", function() {
    fetch("/api/new_game")
        .then(response => response.json())
        .then(data => {
            renderGame(data)
        })
})

newCardBtn.addEventListener("click", function() {
    fetch("/api/hit")
        .then(response => response.json())
        .then(data => {
            renderGame(data)
        })
})
