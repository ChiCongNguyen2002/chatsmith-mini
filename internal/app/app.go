package app

import (
	"context"
	"os/signal"
	"syscall"

	"chatsmith-mini/configs"
	"chatsmith-mini/internal/controller"
	"chatsmith-mini/internal/router"
	"chatsmith-mini/internal/server"
	"chatsmith-mini/pkg/databases"
	"chatsmith-mini/pkg/logger"
)

func Run(cfg *configs.Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := logger.New()

	db, err := databases.NewPostgres(cfg.Database)
	if err != nil {
		return err
	}

	controllers := &controller.Controllers{HealthCheck: controller.NewHealthCheck(db)}
	engine := router.New(log, controllers)

	return server.New(cfg.Server, log, engine).Start(ctx)
}
