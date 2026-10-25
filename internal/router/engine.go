package router

import (
	"log/slog"

	"chatsmith-mini/internal/controller"
	"chatsmith-mini/internal/middleware"
	"chatsmith-mini/internal/repository"
	"chatsmith-mini/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool, log *slog.Logger) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.RequestID(), middleware.AccessLog(log), middleware.Recovery(log))

	healthCheck := controller.NewHealthCheck(db, log)
	engine.GET("/healthz", healthCheck.Check)

	conversation := controller.NewConversation(service.NewConversation(repository.NewConversation(db)), log)
	v1 := engine.Group("/v1", middleware.UserID())
	v1.POST("/conversations", conversation.Create)
	v1.GET("/conversations", conversation.List)
	v1.GET("/conversations/:id/messages", conversation.Messages)

	return engine
}
