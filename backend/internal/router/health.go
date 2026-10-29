package router

import (
	"chatsmith-mini/internal/controller"

	"github.com/gin-gonic/gin"
)

const routeHealthz = "/healthz"

func registerHealth(engine *gin.Engine, h *controller.Health) {
	engine.GET(routeHealthz, h.Check)
}
