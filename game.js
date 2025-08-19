const game = {
    player: {
        name: "Per",
        chips: 200
    },
    cards: [],
    sum: 0,
    hasBlackJack: false,
    isAlive: false,
    message: "",

    getRandomCard() {
        let randomNumber = Math.floor(Math.random() * 13) + 1;
        if (randomNumber > 10) {
            return 10;
        } else if (randomNumber === 1) {
            return 11;
        } else {
            return randomNumber;
        }
    },

    startGame() {
        this.isAlive = true;
        this.hasBlackJack = false;
        let firstCard = this.getRandomCard();
        let secondCard = this.getRandomCard();
        this.cards = [firstCard, secondCard];
        this.sum = firstCard + secondCard;
        return this.updateGameStatus();
    },

    newCard() {
        if (this.isAlive && !this.hasBlackJack) {
            let card = this.getRandomCard();
            this.sum += card;
            this.cards.push(card);
            return this.updateGameStatus();
        }
    },

    updateGameStatus() {
        if (this.sum <= 20) {
            this.message = "Do you want to draw a new card?";
        } else if (this.sum === 21) {
            this.message = "You've got Blackjack!";
            this.hasBlackJack = true;
        } else {
            this.message = "You're out of the game!";
            this.isAlive = false;
        }
        return {
            cards: this.cards,
            sum: this.sum,
            hasBlackJack: this.hasBlackJack,
            isAlive: this.isAlive,
            message: this.message
        };
    },
    
    getGameState() {
        return {
            cards: this.cards, 
            sum: this.sum, 
            hasBlackJack: this.hasBlackJack, 
            isAlive: this.isAlive, 
            message: this.message, 
            player: this.player
        }
    }
};

module.exports = game;