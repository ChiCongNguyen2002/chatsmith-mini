package router

import (
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/controller"

	"github.com/gin-gonic/gin"
)

func registerChat(group *gin.RouterGroup, h *controller.Chat) {
	group.POST(constants.PathConversationMessages, h.Send)
	group.POST(constants.PathConversationContinue, h.Continue)
}

func registerChatMessages(group *gin.RouterGroup, h *controller.Chat) {
	group.GET(constants.PathConversationMessages, h.Messages)
}
