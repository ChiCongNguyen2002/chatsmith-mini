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

type HealthCheck struct {
	db databases.PostgresDB
}

func NewHealthCheck(db databases.PostgresDB) *HealthCheck {
	return &HealthCheck{db: db}
}

func (h *HealthCheck) Check(ctx *gin.Context) {
	if err := h.ping(ctx.Request.Context()); err != nil {
		logger.FromContext(ctx.Request.Context()).Error(logPingFailed, constants.LogKeyError, err)
		response.Err(ctx, failure.CodeServiceUnavailable)
		return
	}

	response.Success(ctx, healthStatusOK)
}

func (h *HealthCheck) ping(ctx context.Context) error {
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
