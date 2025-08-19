# Build stage
FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /blackjack

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /blackjack .
COPY --from=builder /app/frontend ./frontend

EXPOSE 8080

CMD ["./blackjack"]
