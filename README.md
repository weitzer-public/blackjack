# Blackjack

This is a simple web-based Blackjack game with a Go backend and a plain HTML/CSS/JS frontend.

## How to Run Locally

### Prerequisites

- [Go](https://golang.org/doc/install) installed on your machine.

### Steps

1. **Clone the repository:**

   ```bash
   git clone https://github.com/weitzer-public/blackjack.git
   cd blackjack
   ```

2. **Run the backend server:**

   ```bash
   go run .
   ```

3. **Open the application in your browser:**

   Navigate to [http://localhost:8080](http://localhost:8080) in your web browser.

## How to Build with Docker

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) installed on your machine.

### Steps

1. **Build the Docker image:**

   ```bash
   docker build -t blackjack .
   ```

2. **Run the Docker container:**

   ```bash
   docker run -p 8080:8080 blackjack
   ```

3. **Open the application in your browser:**

   Navigate to [http://localhost:8080](http://localhost:8080) in your web browser.
