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
	project.DELETE("/:id/chapters/:chapter_internal_id", projectDeleteChapterAction.Invoke)

	project.POST("/:id/characters", projectCreateCharacterAction.Invoke)
	project.POST("/:id/characters/:character_internal_id", projectUpdateCharacterAction.Invoke)
	project.DELETE("/:id/characters/:character_internal_id", projectDeleteCharacterAction.Invoke)

	project.POST("/:id/locations", projectCreateLocationAction.Invoke)
	project.POST("/:id/locations/:location_internal_id", projectUpdateLocationAction.Invoke)
	project.DELETE("/:id/locations/:location_internal_id", projectDeleteLocationAction.Invoke)
}
