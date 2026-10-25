package middleware

import (
	"net/http"
	"strings"

	"chatsmith-mini/pkg/response"
	"chatsmith-mini/pkg/uid"

	"github.com/gin-gonic/gin"
)

const (
	UserIDHeader = "X-User-ID"
	userIDKey    = "user_id"
)

func UserID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := strings.ToLower(ctx.GetHeader(UserIDHeader))
		if !uid.Valid(id) {
			response.Error(ctx, http.StatusUnauthorized, "unauthorized")
			ctx.Abort()
			return
		}

		ctx.Set(userIDKey, id)
		ctx.Next()
	}
}

func GetUserID(ctx *gin.Context) string {
	return ctx.GetString(userIDKey)
}
