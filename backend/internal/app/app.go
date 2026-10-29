package app

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/auth"
	"chatsmith-mini/internal/catalog"
	"chatsmith-mini/internal/chat"
	"chatsmith-mini/internal/controller"
	"chatsmith-mini/internal/conversation"
	"chatsmith-mini/internal/middleware"
	"chatsmith-mini/internal/provider"
	"chatsmith-mini/internal/router"
	"chatsmith-mini/internal/server"
	"chatsmith-mini/internal/user"
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

	httpClient := provider.NewHTTPClient()
	streamers := newStreamers(cfg.Provider, httpClient)
	available := func(p constants.Provider) bool { _, ok := streamers[p]; return ok }

	users := user.NewModule(&user.CreateModuleDTO{DB: db})
	authModule := auth.NewModule(&auth.CreateModuleDTO{Config: cfg.Auth, Users: users.Service})
	conversations := conversation.NewModule(&conversation.CreateModuleDTO{DB: db})
	catalogs := catalog.NewModule(&catalog.CreateModuleDTO{DB: db, Available: available, DefaultModel: cfg.Chat.DefaultModel})
	if err := catalogs.Service.CheckDefault(ctx); err != nil {
		return err
	}
	classifier := provider.NewClassifier(cfg.Provider, cfg.Chat, httpClient)
	chats := chat.NewModule(&chat.CreateModuleDTO{
		DB:            db,
		Catalog:       catalogs.Service,
		Conversations: conversations.Service,
		Streamers:     streamers,
		Classifier:    classifier,
		Config:        cfg.Chat,
	})

	controllers := &controller.Controllers{
		Web:          controller.NewWeb(cfg.Server.FrontendDir),
		HealthCheck:  controller.NewHealthCheck(db),
		Conversation: controller.NewConversation(conversations.Service),
		Catalog:      controller.NewCatalog(catalogs.Service),
		Chat:         controller.NewChat(chats.Service, cfg.Chat.MaxBodyBytes),
		Auth:         controller.NewAuth(authModule.Service, cfg.Auth),
	}
	engine := router.New(log, controllers, cfg.Server.RequestTimeout, cfg.Server.FrontendDir, middleware.Session(authModule.Service))

	return server.New(cfg.Server, log, engine).Start(ctx)
}

func newStreamers(cfg configs.ProviderConfig, client *http.Client) map[constants.Provider]chat.Streamer {
	streamers := make(map[constants.Provider]chat.Streamer)
	if cfg.OpenAIAPIKey != "" {
		streamers[constants.ProviderOpenAI] = provider.NewOpenAI(cfg, client)
	}
	if cfg.AnthropicAPIKey != "" {
		streamers[constants.ProviderAnthropic] = provider.NewAnthropic(cfg, client)
	}

	return streamers
}
