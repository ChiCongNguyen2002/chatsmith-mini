package router

import (
	"chatsmith-mini/internal/controller"

	"github.com/gin-gonic/gin"
)

const routeModels = "/models"

func registerCatalog(group *gin.RouterGroup, h *controller.Catalog) {
	group.GET(routeModels, h.Models)
}
