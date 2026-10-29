package controller

import (
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

func respondErr(ctx *gin.Context, err error) {
	response.ErrFrom(ctx, err)
}
