package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	RequestIDHeader = "X-Request-ID"
	requestIDKey    = "request_id"
	maxRequestIDLen = 64
)

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetHeader(RequestIDHeader)
		if id == "" || len(id) > maxRequestIDLen {
			id = newRequestID()
		}

		ctx.Set(requestIDKey, id)
		ctx.Header(RequestIDHeader, id)
		ctx.Next()
	}
}

func AccessLog(log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		log.Info("http_request",
			"request_id", ctx.GetString(requestIDKey),
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"status", ctx.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
		)
	}
}

func Recovery(log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error("panic_recovered",
					"request_id", ctx.GetString(requestIDKey),
					"panic", recovered,
					"stack", string(debug.Stack()),
				)
				response.Error(ctx, http.StatusInternalServerError, "internal error")
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

func newRequestID() string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)

	return hex.EncodeToString(buf)
}
