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

    sumEl.textContent = "Sum: " + data.PlayerScore

    if (data.PlayerScore <= 20) {
        messageEl.textContent = "Do you want to draw a new card?"
    } else if (data.PlayerScore === 21) {
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
