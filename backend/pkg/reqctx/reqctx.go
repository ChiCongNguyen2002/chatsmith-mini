package reqctx

import (
	"chatsmith-mini/constants"

	"github.com/gin-gonic/gin"
)

func UserID(ctx *gin.Context) string {
	return ctx.GetString(string(constants.ContextKeyUserID))
}
