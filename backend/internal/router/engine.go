package router

import (
	"log/slog"
	"time"

	"chatsmith-mini/internal/auth"
	"chatsmith-mini/internal/controller"
	"chatsmith-mini/internal/middleware"

	"github.com/gin-gonic/gin"
)

const (
	routeRoot             = "/"
	routeConversationPage = "/conversation/:" + controller.ParamConversationID
	routeHealthz          = "/healthz"
	routeV1               = "/v1"
	routeStatic           = "/static"
)

func New(log *slog.Logger, controllers *controller.Controllers, requestTimeout time.Duration, frontendDir string, session gin.HandlerFunc) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.RequestID(log), middleware.AccessLog(), middleware.Recovery())

	engine.GET(routeRoot, controllers.Web.Index)
	engine.GET(routeConversationPage, controllers.Web.Index)
	engine.Static(routeStatic, frontendDir)
	engine.GET(routeHealthz, controllers.HealthCheck.Check)
	engine.GET(auth.PathGoogleLogin, controllers.Auth.Login)
	engine.GET(auth.PathGoogleCallback, controllers.Auth.Callback)
	engine.POST(auth.PathLogout, controllers.Auth.Logout)

	v1 := engine.Group(routeV1, session)
	registerAuth(v1, controllers.Auth)
	registerChat(v1, controllers.Chat)

	timed := v1.Group("", middleware.Timeout(requestTimeout))
	registerConversation(timed, controllers.Conversation)
	registerCatalog(timed, controllers.Catalog)

	return engine
}
