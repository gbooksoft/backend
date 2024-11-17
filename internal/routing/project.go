package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func projectRoutes(e *gin.Engine) {
	project := e.Group("/projects")

	project.Use(middleware.AuthRequired())

	project.POST("/", actionProject.Create)
	project.GET("/:id/menu", actionProject.GetProjectMenu)
	project.POST("/:id/menu", actionProject.AddMenuItem)
	project.POST("/:id/menu/:item_id/toggle-expand", actionProject.ToggleExpandMenuItem)
}
