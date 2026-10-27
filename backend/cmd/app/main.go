package main

import (
	"context"
	"log"

	"chatsmith-mini/configs"
	"chatsmith-mini/internal/app"
)

func main() {
	cfg, err := configs.Load(context.Background())
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatalf("app exit: %v", err)
	}
}
