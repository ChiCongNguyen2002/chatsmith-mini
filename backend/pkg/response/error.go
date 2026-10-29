package response

import (
	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/logger"

	"github.com/gin-gonic/gin"
)

const logRequestFailed = "request_failed"

func ErrFrom(ctx *gin.Context, err error) {
	code := failure.CodeOf(err)
	if code == failure.CodeInternal {
		logger.FromContext(ctx.Request.Context()).Error(logRequestFailed, constants.LogKeyError, err)
	}

	Err(ctx, code)
}
