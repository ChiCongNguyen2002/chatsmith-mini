package middleware

import (
	"context"

	"chatsmith-mini/constants"
	"chatsmith-mini/internal/domain/auth"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

type SessionVerifier interface {
	VerifySession(ctx context.Context, token string) (string, error)
}

func Session(sessions SessionVerifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie(auth.CookieSession)
		if err != nil {
			response.Err(ctx, failure.CodeUnauthorized)
			ctx.Abort()
			return
		}

		userID, err := sessions.VerifySession(ctx.Request.Context(), token)
		if err != nil {
			response.Err(ctx, failure.CodeUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Set(string(constants.ContextKeyUserID), userID)
		ctx.Next()
	}
}
