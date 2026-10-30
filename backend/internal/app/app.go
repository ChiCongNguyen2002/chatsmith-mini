package app

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/controller"
	"chatsmith-mini/internal/domain/auth"
	"chatsmith-mini/internal/domain/catalog"
	"chatsmith-mini/internal/domain/chat"
	"chatsmith-mini/internal/domain/conversation"
	"chatsmith-mini/internal/domain/user"
	"chatsmith-mini/internal/middleware"
	"chatsmith-mini/internal/router"
	"chatsmith-mini/internal/server"
	"chatsmith-mini/pkg/databases"
	"chatsmith-mini/pkg/llm"
	"chatsmith-mini/pkg/logger"
	"chatsmith-mini/pkg/telemetry"
)

const logTelemetryShutdown = "telemetry_shutdown"

func Run(cfg *configs.Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := logger.New()

	shutdownMetrics, err := telemetry.Setup(ctx, cfg.Telemetry.ProjectID)
	if err != nil {
		return err
	}
	defer func() {
		if err := shutdownMetrics(context.Background()); err != nil {
			log.Error(logTelemetryShutdown, constants.LogKeyError, err)
		}
	}()

	db, err := databases.NewPostgres(cfg.Database)
	if err != nil {
		return err
	}

	httpClient := llm.NewHTTPClient()
	streamers := newStreamers(cfg.Provider, httpClient)
	available := func(p constants.Provider) bool { _, ok := streamers[p]; return ok }

	users := user.NewModule(db)
	authModule := auth.NewModule(db, users.Service, newAuthProviders(cfg.Auth), cfg.Auth)
	conversations := conversation.NewModule(db)
	catalogs := catalog.NewModule(db, available, cfg.Chat.DefaultModel)
	if err := catalogs.Service.CheckDefault(ctx); err != nil {
		return err
	}
	chats, err := chat.NewModule(db, conversations.Service, catalogs.Service, streamers, llm.NewClassifier(cfg.Provider, cfg.Chat, httpClient), cfg.Chat)
	if err != nil {
		return err
	}

	engine, err := router.NewEngine(log, cfg.Server.RequestTimeout, middleware.Session(authModule.Service), router.Handlers{
		Web:          controller.NewWeb(cfg.Server.FrontendDir),
		Health:       controller.NewHealth(db),
		Auth:         controller.NewAuth(authModule.Service, cfg.Auth),
		Catalog:      controller.NewCatalog(catalogs.Service),
		Conversation: controller.NewConversation(conversations.Service),
		Chat:         controller.NewChat(chats.Service, cfg.Chat.MaxBodyBytes),
	})
	if err != nil {
		return err
	}

	return server.New(cfg.Server, log, engine).Start(ctx)
}

func newStreamers(cfg configs.ProviderConfig, client *http.Client) map[constants.Provider]chat.Streamer {
	streamers := make(map[constants.Provider]chat.Streamer)
	if cfg.OpenAIAPIKey != "" {
		streamers[constants.ProviderOpenAI] = llm.NewOpenAI(cfg, client)
	}
	if cfg.AnthropicAPIKey != "" {
		streamers[constants.ProviderAnthropic] = llm.NewAnthropic(cfg, client)
	}

	return streamers
}
