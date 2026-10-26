package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func Timeout(limit time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, cancel := context.WithTimeout(ctx.Request.Context(), limit)
		defer cancel()

		ctx.Request = ctx.Request.WithContext(limited)
		ctx.Next()
	}
}
