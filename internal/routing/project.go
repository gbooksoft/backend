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
}
