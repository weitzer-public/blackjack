I have addressed the feedback from the code review:

1.  **Refactored Frontend `fetch` Calls:**
    *   I created a new helper function called `performAction(url, method = 'GET', data = null)` in `frontend/script.js`.
    *   This function takes the API endpoint URL, HTTP method (defaults to 'GET'), and optional data for POST requests.
    *   It handles the `fetch` call, including setting the appropriate headers and body for POST requests.
    *   All the event listeners for the buttons (`newGameBtn`, `betBtn`, `hitBtn`, `standBtn`, `doubleDownBtn`, `splitBtn`) are now using this new function, which has significantly reduced code duplication and improved maintainability.

2.  **Added Frontend Error Handling:**
    *   The `performAction` function includes robust error handling for the `fetch` calls.
    *   It checks if the response from the server is `ok`. If not, it throws an error with the error message from the server's JSON response.
    *   The `catch` block now catches these errors and displays a user-friendly error message in the `messageBar` on the UI.
    *   This ensures that if the server returns an error, the user is notified about the problem, making the application more robust.