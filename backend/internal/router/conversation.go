package router

import (
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/controller"

	"github.com/gin-gonic/gin"
)

func registerConversation(group *gin.RouterGroup, h *controller.Conversation) {
	group.POST(constants.PathConversations, h.Create)
	group.GET(constants.PathConversations, h.List)
	group.GET(constants.PathConversation, h.Get)
	group.DELETE(constants.PathConversation, h.Delete)
	group.PATCH(constants.PathConversation, h.Rename)
}
