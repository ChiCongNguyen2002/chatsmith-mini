package controller

import (
	"net/http"

	"chatsmith-mini/constants"
	"chatsmith-mini/internal/domain/chat"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/reqctx"
	"chatsmith-mini/pkg/response"
	"chatsmith-mini/pkg/sse"

	"github.com/gin-gonic/gin"
)

type Chat struct {
	service      *chat.Service
	maxBodyBytes int64
}

func NewChat(service *chat.Service, maxBodyBytes int64) *Chat {
	return &Chat{service: service, maxBodyBytes: maxBodyBytes}
}

func (c *Chat) Send(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, c.maxBodyBytes)

	var req chat.SendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondErr(ctx, failure.New(failure.CodeBadRequest))
		return
	}

	turn, err := c.service.Start(ctx.Request.Context(), chat.StartInput{
		UserID:         reqctx.UserID(ctx),
		ConversationID: ctx.Param(constants.ParamConversationID),
		Content:        req.Content,
		Model:          req.Model,
		RequestID:      ctx.GetString(string(constants.ContextKeyRequestID)),
	})
	if err != nil {
		respondErr(ctx, err)
		return
	}

	c.run(ctx, turn)
}

func (c *Chat) Continue(ctx *gin.Context) {
	turn, err := c.service.Continue(ctx.Request.Context(), chat.ContinueInput{
		UserID:         reqctx.UserID(ctx),
		ConversationID: ctx.Param(constants.ParamConversationID),
		RequestID:      ctx.GetString(string(constants.ContextKeyRequestID)),
	})
	if err != nil {
		respondErr(ctx, err)
		return
	}

	c.run(ctx, turn)
}

func (c *Chat) Messages(ctx *gin.Context) {
	var query chat.MessageListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		respondErr(ctx, failure.New(failure.CodeBadRequest))
		return
	}

	items, err := c.service.ListMessages(ctx.Request.Context(), reqctx.UserID(ctx), ctx.Param(constants.ParamConversationID), query.AfterSeq, query.Limit)
	if err != nil {
		respondErr(ctx, err)
		return
	}

	response.Success(ctx, chat.ToMessageResponses(items))
}

func (c *Chat) run(ctx *gin.Context, turn *chat.Turn) {
	stream := sse.Open(ctx)
	defer stream.Close()

	c.service.Run(ctx.Request.Context(), turn, stream.Send)
}
