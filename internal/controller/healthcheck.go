package controller

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"chatsmith-mini/internal/middleware"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pingTimeout = 2 * time.Second

type HealthCheck struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewHealthCheck(db *pgxpool.Pool, log *slog.Logger) *HealthCheck {
	return &HealthCheck{db: db, log: log}
}

func (h *HealthCheck) Check(ctx *gin.Context) {
	pingCtx, cancel := context.WithTimeout(ctx.Request.Context(), pingTimeout)
	defer cancel()

	if err := h.db.Ping(pingCtx); err != nil {
		h.log.Error("healthcheck_db_ping",
			"request_id", ctx.Writer.Header().Get(middleware.RequestIDHeader),
			"error", failure.ErrWithTrace(err),
		)
		response.Error(ctx, http.StatusServiceUnavailable, "database unavailable")
		return
	}

	response.Success(ctx, "ok")
}
