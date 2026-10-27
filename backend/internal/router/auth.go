package router

import (
	"chatsmith-mini/internal/controller"

	"github.com/gin-gonic/gin"
)

const routeMe = "/me"

func registerAuth(group *gin.RouterGroup, c *controller.Auth) {
	group.GET(routeMe, c.Me)
}
