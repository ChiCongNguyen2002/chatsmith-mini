package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPResponse struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, HTTPResponse{Status: http.StatusOK, Data: data})
}

func Error(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, HTTPResponse{Status: status, Data: message})
}
