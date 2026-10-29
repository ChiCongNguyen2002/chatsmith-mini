.PHONY: install build test lint

install:
	go mod download

build:
	go build -o bin/app ./cmd/app

test:
	go test ./...

lint:
	golangci-lint run ./...
