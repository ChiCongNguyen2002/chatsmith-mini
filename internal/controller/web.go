package controller

import (
	"net/http"

	"chatsmith-mini/web"

	"github.com/gin-gonic/gin"
)

const contentTypeHTML = "text/html; charset=utf-8"

type Web struct{}

func NewWeb() *Web {
	return &Web{}
}

func (w *Web) Index(ctx *gin.Context) {
	ctx.Data(http.StatusOK, contentTypeHTML, web.IndexHTML)
}
