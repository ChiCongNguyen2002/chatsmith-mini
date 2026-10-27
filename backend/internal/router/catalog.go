package router

import (
	"chatsmith-mini/internal/controller"

	"github.com/gin-gonic/gin"
)

const routeModels = "/models"

func registerCatalog(group *gin.RouterGroup, c *controller.Catalog) {
	group.GET(routeModels, c.Models)
}
