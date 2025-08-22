I have addressed all the issues discovered by the code review. Here are the details on how I fixed each issue:

*   **`main.go`**:
    *   Introduced a `sync.Mutex` to protect the global `game` variable from concurrent access, fixing a critical race condition.
    *   Updated the `hitHandler` to use the `Bust` constant from the `PlayerStatus` enum.
*   **`game.go`**:
    *   Replaced a hardcoded `17` with the `DealerStand` constant in the computer player's and dealer's logic to improve maintainability.
    *   Introduced a new `PlayerStatus` type with an `iota`-based enumeration to replace the string-based `Status` field, improving type safety.
    *   Updated the `Player` struct and all related logic to use the new `PlayerStatus` type.
*   **`frontend/script.js`**:
    *   Improved the user experience by providing more specific "Game over!" messages based on the human player's final status.
    *   Removed the unreachable `default` case from the game state `switch` statement.