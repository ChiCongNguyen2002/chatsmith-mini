package controller

import (
	"chatsmith-mini/internal/catalog"
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

type Catalog struct {
	service *catalog.Service
}

func NewCatalog(service *catalog.Service) *Catalog {
	return &Catalog{service: service}
}

func (c *Catalog) Models(ctx *gin.Context) {
	result, err := c.service.List(ctx.Request.Context())
	if err != nil {
		respondError(ctx, err)
		return
	}

	response.Success(ctx, result)
}
