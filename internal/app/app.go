package app

import (
	"context"
	"os/signal"
	"syscall"

	"chatsmith-mini/configs"
	"chatsmith-mini/internal/controller"
	"chatsmith-mini/internal/conversation"
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

	conversations := conversation.NewModule(&conversation.CreateModuleDTO{DB: db})
	controllers := &controller.Controllers{
		Web:          controller.NewWeb(),
		HealthCheck:  controller.NewHealthCheck(db),
		Conversation: controller.NewConversation(conversations.Service),
	}
	engine := router.New(log, controllers, cfg.Server.RequestTimeout)

	return server.New(cfg.Server, log, engine).Start(ctx)
}
