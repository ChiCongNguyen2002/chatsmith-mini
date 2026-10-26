package middleware

import (
	"time"

	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/logger"

	"github.com/gin-gonic/gin"
)

const logHTTPRequest = "http_request"

func AccessLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		logger.FromContext(ctx.Request.Context()).Info(logHTTPRequest,
			constants.LogKeyMethod, ctx.Request.Method,
			constants.LogKeyPath, ctx.Request.URL.Path,
			constants.LogKeyStatus, ctx.Writer.Status(),
			constants.LogKeyLatencyMS, time.Since(start).Milliseconds(),
		)
	}
}
