### Plan

- [ ] **Refactor tests for determinism:**
    - In `game_test.go`, I will create a new function `newTestGame()` that initializes a `Game` with a predictable, unshuffled deck. This will ensure our tests are reliable and not subject to random shuffling.
    - I will then update `TestNewGame`, `TestHit`, and `TestStand` to use this new `newTestGame()` function.
- [ ] **Correct game logic in `game.go`:**
    - I will modify the `NewGame()` function to deal two cards to the dealer from the start, which is the standard in Blackjack. I'll keep the second card concealed until the appropriate time.
    - I will update `NewGame()` to accurately check for a natural blackjack (21) for both the player and the dealer, and I'll set the game state to "player_blackjack", "dealer_blackjack", or "push" as needed.
    - I will adjust the `Visible()` function to reveal the dealer's complete hand and score only when the game has concluded.
- [ ] **Verify the fix:**
    - I will run the updated tests to confirm that they pass and that the game logic is sound.

I will now proceed with this plan. I will start by refactoring the tests for determinism.