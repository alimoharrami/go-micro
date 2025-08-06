# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Install git (needed for go mod download with private/public repos)
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd

# Stage 2: Run the Go binary in a minimal image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

CMD ["./main"]
