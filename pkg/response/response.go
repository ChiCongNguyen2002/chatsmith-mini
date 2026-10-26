package response

import (
	"net/http"

	"chatsmith-mini/pkg/failure"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Body{Status: http.StatusOK, Data: data})
}

func Err(ctx *gin.Context, code failure.Code) {
	status := code.HTTPStatus()
	ctx.JSON(status, Body{Status: status, Data: string(code)})
}
