package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
}

func project2Routes(e *gin.Engine) {
	project := e.Group("/projects2")

	project.Use(middleware.AuthRequired())

	project.POST("", projectCreateAction.Invoke)
	project.POST("/:id", projectUpdateAction.Invoke)
	project.POST("/:id/chapters", projectCreateChapterAction.Invoke)
	project.POST("/:id/chapters/:chapter_internal_id", projectUpdateChapterAction.Invoke)
}
