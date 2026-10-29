package controller

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const (
	fileIndex = "index.html"
)

type Web struct {
	indexPath   string
	frontendDir string
}

func NewWeb(frontendDir string) *Web {
	return &Web{indexPath: filepath.Join(frontendDir, fileIndex), frontendDir: frontendDir}
}

func (w *Web) Index(ctx *gin.Context) {
	ctx.File(w.indexPath)
}

func (w *Web) StaticDir() string {
	return w.frontendDir
}
