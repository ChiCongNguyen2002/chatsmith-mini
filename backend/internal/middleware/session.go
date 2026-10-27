package middleware

import (
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/auth"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

func Session(sessions *auth.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie(auth.CookieSession)
		if err != nil {
			response.Err(ctx, failure.CodeUnauthorized)
			ctx.Abort()
			return
		}

		userID, err := sessions.VerifySession(token)
		if err != nil {
			response.Err(ctx, failure.CodeUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Set(string(constants.ContextKeyUserID), userID)
		ctx.Next()
	}
}

func GetUserID(ctx *gin.Context) string {
	return ctx.GetString(string(constants.ContextKeyUserID))
}
