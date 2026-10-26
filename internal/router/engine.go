package router

import (
	"log/slog"
	"time"

	"chatsmith-mini/internal/controller"
	"chatsmith-mini/internal/middleware"

	"github.com/gin-gonic/gin"
)

const (
	routeRoot             = "/"
	routeConversationPage = "/conversation/:" + controller.ParamConversationID
	routeHealthz          = "/healthz"
	routeV1               = "/v1"
)

func New(log *slog.Logger, controllers *controller.Controllers, requestTimeout time.Duration) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.RequestID(log), middleware.AccessLog(), middleware.Recovery())

	engine.GET(routeRoot, controllers.Web.Index)
	engine.GET(routeConversationPage, controllers.Web.Index)
	engine.GET(routeHealthz, controllers.HealthCheck.Check)

	v1 := engine.Group(routeV1, middleware.UserID())
	registerChat(v1, controllers.Chat)

	timed := v1.Group("", middleware.Timeout(requestTimeout))
	registerConversation(timed, controllers.Conversation)
	registerCatalog(timed, controllers.Catalog)

	return engine
}
