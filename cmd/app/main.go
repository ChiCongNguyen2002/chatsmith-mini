package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chatsmith-mini/configs"
	"chatsmith-mini/internal/router"
	"chatsmith-mini/pkg/databases"
	"chatsmith-mini/pkg/logger"
)

const shutdownTimeout = 8 * time.Second

func main() {
	log := logger.New()

	if err := run(log); err != nil {
		log.Error("app_exit", "error", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := configs.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := databases.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := databases.Migrate(ctx, db, cfg.ScriptsDir); err != nil {
		return err
	}

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router.New(db, log),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Info("server_start", "port", cfg.Port)

	serverErr := make(chan error, 1)
	go func() { serverErr <- server.ListenAndServe() }()

	select {
	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	log.Info("server_shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}
