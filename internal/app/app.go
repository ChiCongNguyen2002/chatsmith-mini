package app

import (
	"context"
	"os/signal"
	"syscall"

	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/catalog"
	"chatsmith-mini/internal/chat"
	"chatsmith-mini/internal/controller"
	"chatsmith-mini/internal/conversation"
	"chatsmith-mini/internal/provider"
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

	streamers := newStreamers(cfg.Provider)
	available := func(p constants.Provider) bool { _, ok := streamers[p]; return ok }

	conversations := conversation.NewModule(&conversation.CreateModuleDTO{DB: db})
	catalogs := catalog.NewModule(&catalog.CreateModuleDTO{DB: db, Available: available, DefaultModel: cfg.Chat.DefaultModel})
	if err := catalogs.Service.CheckDefault(ctx); err != nil {
		return err
	}
	chats := chat.NewModule(&chat.CreateModuleDTO{
		DB:            db,
		Catalog:       catalogs.Service,
		Conversations: conversations.Service,
		Streamers:     streamers,
		Config:        cfg.Chat,
	})

	controllers := &controller.Controllers{
		Web:          controller.NewWeb(),
		HealthCheck:  controller.NewHealthCheck(db),
		Conversation: controller.NewConversation(conversations.Service),
		Catalog:      controller.NewCatalog(catalogs.Service),
		Chat:         controller.NewChat(chats.Service, cfg.Chat.MaxBodyBytes),
	}
	engine := router.New(log, controllers, cfg.Server.RequestTimeout)

	return server.New(cfg.Server, log, engine).Start(ctx)
}

func newStreamers(cfg configs.ProviderConfig) map[constants.Provider]chat.Streamer {
	client := provider.NewHTTPClient()
	streamers := make(map[constants.Provider]chat.Streamer)
	if cfg.OpenAIAPIKey != "" {
		streamers[constants.ProviderOpenAI] = provider.NewOpenAI(cfg, client)
	}
	if cfg.AnthropicAPIKey != "" {
		streamers[constants.ProviderAnthropic] = provider.NewAnthropic(cfg, client)
	}

	return streamers
}
