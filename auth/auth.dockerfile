# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Install git and ssh client
RUN apk add --no-cache git openssh

WORKDIR /app


# Set GOPRIVATE
ENV GOPRIVATE=github.com/alimoharrami/go-micro*

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
