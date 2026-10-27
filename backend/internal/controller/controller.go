package controller

import (
	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/logger"
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

const logRequestFailed = "request_failed"

type Controllers struct {
	Web          *Web
	HealthCheck  *HealthCheck
	Conversation *Conversation
	Catalog      *Catalog
	Chat         *Chat
	Auth         *Auth
}

func respondError(ctx *gin.Context, err error) {
	code := failure.CodeOf(err)
	if code == failure.CodeInternal {
		logger.FromContext(ctx.Request.Context()).Error(logRequestFailed, constants.LogKeyError, err)
	}

	response.Err(ctx, code)
}
