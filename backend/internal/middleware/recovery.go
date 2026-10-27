package middleware

import (
	"runtime/debug"

	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/logger"
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

const logPanicRecovered = "panic_recovered"

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.FromContext(ctx.Request.Context()).Error(logPanicRecovered,
					constants.LogKeyPanic, recovered,
					constants.LogKeyStack, string(debug.Stack()),
				)
				response.Err(ctx, failure.CodeInternal)
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
