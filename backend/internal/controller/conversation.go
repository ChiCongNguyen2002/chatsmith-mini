package controller

import (
	"errors"
	"io"

	"chatsmith-mini/constants"
	"chatsmith-mini/internal/domain/conversation"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/reqctx"
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

type Conversation struct {
	service *conversation.Service
}

func NewConversation(service *conversation.Service) *Conversation {
	return &Conversation{service: service}
}

func (c *Conversation) Create(ctx *gin.Context) {
	var req conversation.CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		respondErr(ctx, failure.New(failure.CodeBadRequest))
		return
	}

	item, err := c.service.Create(ctx.Request.Context(), reqctx.UserID(ctx), req.Title)
	if err != nil {
		respondErr(ctx, err)
		return
	}

	response.Created(ctx, conversation.ToResponse(*item))
}

func (c *Conversation) Get(ctx *gin.Context) {
	item, err := c.service.Get(ctx.Request.Context(), reqctx.UserID(ctx), ctx.Param(constants.ParamConversationID))
	if err != nil {
		respondErr(ctx, err)
		return
	}

	response.Success(ctx, conversation.ToResponse(*item))
}

func (c *Conversation) List(ctx *gin.Context) {
	var query conversation.ListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		respondErr(ctx, failure.New(failure.CodeBadRequest))
		return
	}

	page, err := c.service.List(ctx.Request.Context(), reqctx.UserID(ctx), query.Cursor, query.Limit)
	if err != nil {
		respondErr(ctx, err)
		return
	}

	response.Paged(ctx, conversation.ToResponses(page.Items), page.NextCursor)
}

func (c *Conversation) Delete(ctx *gin.Context) {
	err := c.service.Delete(ctx.Request.Context(), reqctx.UserID(ctx), ctx.Param(constants.ParamConversationID))
	if err != nil {
		respondErr(ctx, err)
		return
	}

	response.NoContent(ctx)
}

func (c *Conversation) Rename(ctx *gin.Context) {
	var req conversation.RenameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondErr(ctx, failure.New(failure.CodeBadRequest))
		return
	}

	err := c.service.Rename(ctx.Request.Context(), reqctx.UserID(ctx), ctx.Param(constants.ParamConversationID), req.Title)
	if err != nil {
		respondErr(ctx, err)
		return
	}

	response.NoContent(ctx)
}
