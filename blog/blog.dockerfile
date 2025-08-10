# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Install git and ssh client
RUN apk add --no-cache git openssh

WORKDIR /app


RUN git config --global url."https://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@github.com/".insteadOf "https://github.com/"


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
COPY --from=builder /app/web ./web

CMD ["./main"]
