package router

import (
	"chatsmith-mini/internal/controller"

	"github.com/gin-gonic/gin"
)

const (
	routeConversations        = "/conversations"
	routeConversation         = "/conversations/:" + controller.ParamConversationID
	routeConversationMessages = routeConversation + "/messages"
)

func registerConversation(group *gin.RouterGroup, c *controller.Conversation) {
	group.POST(routeConversations, c.Create)
	group.GET(routeConversations, c.List)
	group.GET(routeConversation, c.Get)
	group.DELETE(routeConversation, c.Delete)
	group.GET(routeConversationMessages, c.Messages)
}
