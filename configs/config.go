package configs

import (
	"errors"
	"os"

	"chatsmith-mini/pkg/failure"
)

type Config struct {
	Port        string
	DatabaseURL string
	ScriptsDir  string
}

func Load() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, failure.ErrWithTrace(errors.New("DATABASE_URL is required"))
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: databaseURL,
		ScriptsDir:  getEnv("SCRIPTS_DIR", "scripts"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
