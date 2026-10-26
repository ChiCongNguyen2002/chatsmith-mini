package middleware

import (
	"crypto/rand"
	"log/slog"
	"regexp"

	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/logger"

	"github.com/gin-gonic/gin"
)

var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

func RequestID(log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetHeader(constants.HeaderRequestID)
		if !requestIDPattern.MatchString(id) {
			id = rand.Text()
		}

		ctx.Header(constants.HeaderRequestID, id)
		ctx.Set(string(constants.ContextKeyRequestID), id)

		scoped := log.With(constants.LogKeyRequestID, id)
		ctx.Request = ctx.Request.WithContext(logger.WithContext(ctx.Request.Context(), scoped))
		ctx.Next()
	}
}
