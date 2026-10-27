package controller

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const fileIndex = "index.html"

type Web struct {
	indexPath string
}

func NewWeb(frontendDir string) *Web {
	return &Web{indexPath: filepath.Join(frontendDir, fileIndex)}
}

func (w *Web) Index(ctx *gin.Context) {
	ctx.File(w.indexPath)
}
