package router

import (
	"log/slog"

	"chatsmith-mini/internal/controller"
	"chatsmith-mini/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool, log *slog.Logger) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.RequestID(), middleware.AccessLog(log), middleware.Recovery(log))

	healthCheck := controller.NewHealthCheck(db, log)
	engine.GET("/healthz", healthCheck.Check)

	return engine
}
