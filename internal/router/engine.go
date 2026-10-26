package router

import (
	"log/slog"

	"chatsmith-mini/internal/controller"
	"chatsmith-mini/internal/middleware"

	"github.com/gin-gonic/gin"
)

const routeHealthz = "/healthz"

func New(log *slog.Logger, controllers *controller.Controllers) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.RequestID(log), middleware.AccessLog(), middleware.Recovery())

	engine.GET(routeHealthz, controllers.HealthCheck.Check)

	return engine
}
