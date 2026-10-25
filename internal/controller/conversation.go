package controller

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"chatsmith-mini/internal/middleware"
	"chatsmith-mini/internal/service"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

type Conversation struct {
	service *service.Conversation
	log     *slog.Logger
}

func NewConversation(service *service.Conversation, log *slog.Logger) *Conversation {
	return &Conversation{service: service, log: log}
}

type createConversationRequest struct {
	Title string `json:"title"`
}

func (c *Conversation) Create(ctx *gin.Context) {
	var req createConversationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		response.Error(ctx, http.StatusBadRequest, "invalid body")
		return
	}

	item, err := c.service.Create(ctx.Request.Context(), middleware.GetUserID(ctx), req.Title)
	if err != nil {
		c.fail(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response.HTTPResponse{Status: http.StatusCreated, Data: item})
}

func (c *Conversation) List(ctx *gin.Context) {
	limit, ok := parseQueryInt(ctx, "limit", 1)
	if !ok {
		return
	}

	page, err := c.service.List(ctx.Request.Context(), middleware.GetUserID(ctx), ctx.Query("cursor"), int(limit))
	if err != nil {
		c.fail(ctx, err)
		return
	}

	response.Success(ctx, page)
}

func (c *Conversation) Messages(ctx *gin.Context) {
	afterSeq, ok := parseQueryInt(ctx, "after_seq", 0)
	if !ok {
		return
	}
	limit, ok := parseQueryInt(ctx, "limit", 1)
	if !ok {
		return
	}

	items, err := c.service.Messages(ctx.Request.Context(), middleware.GetUserID(ctx), ctx.Param("id"), afterSeq, int(limit))
	if err != nil {
		c.fail(ctx, err)
		return
	}

	response.Success(ctx, items)
}

func (c *Conversation) fail(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, failure.ErrNotFound):
		response.Error(ctx, http.StatusNotFound, "not found")
	case errors.Is(err, failure.ErrBadRequest), errors.Is(err, failure.ErrInvalidCursor):
		response.Error(ctx, http.StatusBadRequest, "bad request")
	default:
		c.log.Error("request_failed",
			"request_id", ctx.Writer.Header().Get(middleware.RequestIDHeader),
			"path", ctx.Request.URL.Path,
			"error", err,
		)
		response.Error(ctx, http.StatusInternalServerError, "internal error")
	}
}

func parseQueryInt(ctx *gin.Context, name string, min int64) (int64, bool) {
	raw := ctx.Query(name)
	if raw == "" {
		return 0, true
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value < min {
		response.Error(ctx, http.StatusBadRequest, "invalid "+name)
		return 0, false
	}

	return value, true
}
