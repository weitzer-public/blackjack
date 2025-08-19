const game = require('./game');

describe('Blackjack Game Logic', () => {

    describe('getRandomCard', () => {
        test('should return a number between 1 and 11', () => {
            const card = game.getRandomCard();
            expect(card).toBeGreaterThanOrEqual(1);
            expect(card).toBeLessThanOrEqual(11);
        });

        test('should return 10 for face cards (random numbers > 10)', () => {
            const mockMath = Object.create(global.Math);
            mockMath.random = () => 0.9;
            global.Math = mockMath;
            const card = game.getRandomCard();
            expect(card).toBe(10);
        });

        test('should return 11 for an Ace (random number is 1)', () => {
            const mockMath = Object.create(global.Math);
            mockMath.random = () => 0;
            global.Math = mockMath;
            const card = game.getRandomCard();
            expect(card).toBe(11);
        });
    });

    describe('Game Flow', () => {
        let getRandomCardSpy;

        beforeEach(() => {
            getRandomCardSpy = jest.spyOn(game, 'getRandomCard');
        });

        afterEach(() => {
            getRandomCardSpy.mockRestore();
        });

        test('startGame should initialize the game correctly', () => {
            getRandomCardSpy.mockReturnValueOnce(5).mockReturnValueOnce(6);
            const state = game.startGame();
            expect(state.cards).toEqual([5, 6]);
            expect(state.sum).toBe(11);
            expect(state.isAlive).toBe(true);
            expect(state.hasBlackJack).toBe(false);
        });

        test('newCard should add a card to the hand', () => {
            getRandomCardSpy.mockReturnValueOnce(5).mockReturnValueOnce(6);
            game.startGame();
            getRandomCardSpy.mockReturnValueOnce(7);
            const state = game.newCard();
            expect(state.cards).toEqual([5, 6, 7]);
            expect(state.sum).toBe(18);
        });

        test('should declare blackjack on 21', () => {
            getRandomCardSpy.mockReturnValueOnce(10).mockReturnValueOnce(11);
            const state = game.startGame();
            expect(state.hasBlackJack).toBe(true);
            expect(state.message).toBe("You've got Blackjack!");
        });

        test('should end the game if sum is over 21', () => {
            getRandomCardSpy.mockReturnValueOnce(10).mockReturnValueOnce(10);
            game.startGame();
            getRandomCardSpy.mockReturnValueOnce(2);
            const state = game.newCard();
            expect(state.isAlive).toBe(false);
            expect(state.message).toBe("You're out of the game!");
        });
    });
});