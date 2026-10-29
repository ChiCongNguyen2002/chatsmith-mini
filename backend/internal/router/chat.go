package router

import (
	"chatsmith-mini/internal/controller"

	"github.com/gin-gonic/gin"
)

func registerChat(group *gin.RouterGroup, c *controller.Chat) {
	group.POST(routeConversationMessages, c.Send)
	group.POST(routeConversationContinue, c.Continue)
}
