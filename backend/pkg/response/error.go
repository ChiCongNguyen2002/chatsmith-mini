package response

import (
	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/logger"

	"github.com/gin-gonic/gin"
)

const (
	logRequestFailed   = "request_failed"
	logRequestRejected = "request_rejected"
)

func ErrFrom(ctx *gin.Context, err error) {
	code := failure.CodeOf(err)
	log := logger.FromContext(ctx.Request.Context())
	if code == failure.CodeInternal {
		log.Error(logRequestFailed, constants.LogKeyErrorCode, code, constants.LogKeyError, err)
	} else {
		log.Warn(logRequestRejected, constants.LogKeyErrorCode, code)
	}

	Err(ctx, code)
}
