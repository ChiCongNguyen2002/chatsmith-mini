package router

import (
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/controller"

	"github.com/gin-gonic/gin"
)

const routeMe = "/me"

func registerAuthPublic(engine *gin.Engine, h *controller.Auth) {
	engine.GET(constants.PathAuthLogin, h.Login)
	engine.GET(constants.PathAuthCallback, h.Callback)
	engine.POST(constants.PathAuthLogout, h.Logout)
}

func registerAuth(group *gin.RouterGroup, h *controller.Auth) {
	group.GET(routeMe, h.Me)
}
