.PHONY: run build test lint clean

APP_NAME=user-service
CMD_PATH=./cmd/main.go

run:
	go run $(CMD_PATH)

build:
	go build -o $(APP_NAME) $(CMD_PATH)

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -f $(APP_NAME)
