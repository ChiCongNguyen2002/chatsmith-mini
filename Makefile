.PHONY: run up down

run:
	go run ./cmd/app

up:
	docker compose up --build -d

down:
	docker compose down
