package controller

import (
	"context"
	"time"

	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/databases"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/logger"
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	healthStatusOK = "ok"
	pingTimeout    = 2 * time.Second
	logPingFailed  = "healthcheck_ping_failed"
)

type Health struct {
	db databases.PostgresDB
}

func NewHealth(db databases.PostgresDB) *Health {
	return &Health{db: db}
}

func (h *Health) Check(ctx *gin.Context) {
	if err := h.ping(ctx.Request.Context()); err != nil {
		logger.FromContext(ctx.Request.Context()).Error(logPingFailed, constants.LogKeyError, err)
		response.Err(ctx, failure.CodeServiceUnavailable)
		return
	}

	response.Success(ctx, healthStatusOK)
}

func (h *Health) ping(ctx context.Context) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return failure.ErrWithTrace(err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		return failure.ErrWithTrace(err)
	}

	return nil
}
