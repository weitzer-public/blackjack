const messageEl = document.getElementById("message-el")
const sumEl = document.getElementById("sum-el")
const cardsEl = document.getElementById("cards-el")
const playerEl = document.getElementById("player-el")
const dealerCardsEl = document.getElementById("dealer-cards-el")
const dealerSumEl = document.getElementById("dealer-sum-el")
const startGameBtn = document.getElementById("start-game-btn")
const newCardBtn = document.getElementById("new-card-btn")
const standBtn = document.getElementById("stand-btn")

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

    dealerCardsEl.textContent = "Dealer's Cards: "
    if (data.State === "playing") {
        dealerCardsEl.textContent += data.Dealer[0].Value + " ?"
        dealerSumEl.textContent = "Dealer's Sum: ?"
    } else {
        for (let i = 0; i < data.Dealer.length; i++) {
            dealerCardsEl.textContent += data.Dealer[i].Value + " "
        }
        dealerSumEl.textContent = "Dealer's Sum: " + data.DealerScore
    }


    switch (data.State) {
        case "playing":
            messageEl.textContent = "Do you want to draw a new card?"
            break;
        case "player_wins":
            messageEl.textContent = "You win!"
            break;
        case "dealer_wins":
            messageEl.textContent = `Dealer wins with ${data.DealerScore}!`
            break;
        case "player_busts":
            messageEl.textContent = "You're out of the game!"
            break;
        case "tie":
            messageEl.textContent = "It's a tie!"
            break;
        default:
            messageEl.textContent = "Game over."
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

standBtn.addEventListener("click", function() {
    fetch("/api/stand")
        .then(response => response.json())
        .then(data => {
            renderGame(data)
        })
})