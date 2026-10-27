package controller

import (
	"net/http"

	"chatsmith-mini/constants"
	"chatsmith-mini/internal/chat"
	"chatsmith-mini/internal/middleware"
	"chatsmith-mini/pkg/failure"

	"github.com/gin-gonic/gin"
)

type Chat struct {
	service      *chat.Service
	maxBodyBytes int64
}

func NewChat(service *chat.Service, maxBodyBytes int64) *Chat {
	return &Chat{service: service, maxBodyBytes: maxBodyBytes}
}

type sendRequest struct {
	Content string `json:"content"`
	Model   string `json:"model"`
}

func (c *Chat) Send(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, c.maxBodyBytes)

	var req sendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondError(ctx, failure.New(failure.CodeBadRequest))
		return
	}

	turn, err := c.service.Start(ctx.Request.Context(), chat.StartInput{
		UserID:         middleware.GetUserID(ctx),
		ConversationID: ctx.Param(ParamConversationID),
		Content:        req.Content,
		Model:          req.Model,
		RequestID:      ctx.GetString(string(constants.ContextKeyRequestID)),
	})
	if err != nil {
		respondError(ctx, err)
		return
	}

	stream := newEventStream(ctx)
	defer stream.stop()

	c.service.Run(ctx.Request.Context(), turn, stream.emit)
}
