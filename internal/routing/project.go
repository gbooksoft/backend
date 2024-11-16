package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func projectRoutes(e *gin.Engine) {
	project := e.Group("/project")

	project.Use(middleware.AuthRequired())

	project.POST("/", actionProject.Create)
}
