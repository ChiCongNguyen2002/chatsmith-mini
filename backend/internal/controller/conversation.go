package controller

import (
	"errors"
	"io"

	"chatsmith-mini/internal/conversation"
	"chatsmith-mini/internal/middleware"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

const ParamConversationID = "id"

type Conversation struct {
	service *conversation.Service
}

func NewConversation(service *conversation.Service) *Conversation {
	return &Conversation{service: service}
}

type createConversationRequest struct {
	Title string `json:"title"`
}

type listConversationsQuery struct {
	Limit  int    `form:"limit" binding:"min=0"`
	Cursor string `form:"cursor"`
}

type listMessagesQuery struct {
	AfterSeq int64 `form:"after_seq" binding:"min=0"`
	Limit    int   `form:"limit" binding:"min=0"`
}

func (c *Conversation) Create(ctx *gin.Context) {
	var req createConversationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		respondError(ctx, failure.New(failure.CodeBadRequest))
		return
	}

	item, err := c.service.Create(ctx.Request.Context(), middleware.GetUserID(ctx), req.Title)
	if err != nil {
		respondError(ctx, err)
		return
	}

	response.Created(ctx, item)
}

func (c *Conversation) Get(ctx *gin.Context) {
	item, err := c.service.Get(ctx.Request.Context(), middleware.GetUserID(ctx), ctx.Param(ParamConversationID))
	if err != nil {
		respondError(ctx, err)
		return
	}

	response.Success(ctx, item)
}

func (c *Conversation) List(ctx *gin.Context) {
	var query listConversationsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		respondError(ctx, failure.New(failure.CodeBadRequest))
		return
	}

	page, err := c.service.List(ctx.Request.Context(), middleware.GetUserID(ctx), query.Cursor, query.Limit)
	if err != nil {
		respondError(ctx, err)
		return
	}

	response.Success(ctx, page)
}

func (c *Conversation) Delete(ctx *gin.Context) {
	err := c.service.Delete(ctx.Request.Context(), middleware.GetUserID(ctx), ctx.Param(ParamConversationID))
	if err != nil {
		respondError(ctx, err)
		return
	}

	response.NoContent(ctx)
}

type renameConversationRequest struct {
	Title string `json:"title" binding:"required"`
}

func (c *Conversation) Rename(ctx *gin.Context) {
	var req renameConversationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondError(ctx, failure.New(failure.CodeBadRequest))
		return
	}

	err := c.service.Rename(ctx.Request.Context(), middleware.GetUserID(ctx), ctx.Param(ParamConversationID), req.Title)
	if err != nil {
		respondError(ctx, err)
		return
	}

	response.NoContent(ctx)
}

func (c *Conversation) Messages(ctx *gin.Context) {
	var query listMessagesQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		respondError(ctx, failure.New(failure.CodeBadRequest))
		return
	}

	items, err := c.service.Messages(ctx.Request.Context(), middleware.GetUserID(ctx), ctx.Param(ParamConversationID), query.AfterSeq, query.Limit)
	if err != nil {
		respondError(ctx, err)
		return
	}

	response.Success(ctx, items)
}
