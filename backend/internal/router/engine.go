package router

import (
	"log/slog"
	"time"

	"chatsmith-mini/constants"
	"chatsmith-mini/internal/controller"
	"chatsmith-mini/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const routeV1 = "/v1"

type Handlers struct {
	Web          *controller.Web
	Health       *controller.Health
	Auth         *controller.Auth
	Catalog      *controller.Catalog
	Conversation *controller.Conversation
	Chat         *controller.Chat
}

func NewEngine(log *slog.Logger, requestTimeout time.Duration, session gin.HandlerFunc, h Handlers) (*gin.Engine, error) {
	httpMetrics, err := middleware.Metrics()
	if err != nil {
		return nil, err
	}

	engine := gin.New()
	engine.Use(otelgin.Middleware(constants.ServiceName), middleware.RequestID(log), middleware.AccessLog(), httpMetrics, middleware.Recovery())

	registerWeb(engine, h.Web)
	registerHealth(engine, h.Health)
	registerAuthPublic(engine, h.Auth)

	v1 := engine.Group(routeV1, session)
	registerAuth(v1, h.Auth)
	registerChat(v1, h.Chat)

	timed := v1.Group("", middleware.Timeout(requestTimeout))
	registerConversation(timed, h.Conversation)
	registerChatMessages(timed, h.Chat)
	registerCatalog(timed, h.Catalog)

	return engine, nil
}
