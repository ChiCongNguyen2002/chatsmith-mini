.PHONY: run up down test

COMPOSE := docker compose --env-file .env -f deploy/docker-compose.yml

run:
	cd backend && go run ./cmd/app

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down

test:
	cd backend && go vet ./... && go build ./...
