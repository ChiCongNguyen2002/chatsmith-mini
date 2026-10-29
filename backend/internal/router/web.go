package router

import (
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/controller"

	"github.com/gin-gonic/gin"
)

const (
	routeRoot             = "/"
	routeConversationPage = "/conversation/:" + constants.ParamConversationID
	routeStatic           = "/static"
)

func registerWeb(engine *gin.Engine, h *controller.Web) {
	engine.GET(routeRoot, h.Index)
	engine.GET(routeConversationPage, h.Index)
	engine.Static(routeStatic, h.StaticDir())
}
