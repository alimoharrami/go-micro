# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency definitions
COPY go.mod ./
# Download dependencies (if you had go.sum, would copy that too)
# RUN go mod download 

COPY . .

# Build the application
RUN go build -o gateway ./cmd/main.go

# Run Stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/gateway .

# KEY STEP: Copy the static files directory
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./gateway"]
