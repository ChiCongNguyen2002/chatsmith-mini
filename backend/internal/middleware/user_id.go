package middleware

import (
	"strings"

	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/response"
	"chatsmith-mini/pkg/uid"

	"github.com/gin-gonic/gin"
)

func UserID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := strings.ToLower(ctx.GetHeader(constants.HeaderUserID))
		if !uid.Valid(id) {
			response.Err(ctx, failure.CodeUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Set(string(constants.ContextKeyUserID), id)
		ctx.Next()
	}
}

func GetUserID(ctx *gin.Context) string {
	return ctx.GetString(string(constants.ContextKeyUserID))
}
