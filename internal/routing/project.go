package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
}

func projectRoutes(e *gin.Engine) {
	project := e.Group("/projects")

	project.Use(middleware.AuthRequired())

	project.POST("", projectCreateAction.Invoke)
	project.GET("", projectGetAllAction.Invoke)
	project.POST("/:id", projectUpdateAction.Invoke)
	project.GET("/:id", projectGetAction.Invoke)

	project.POST("/:id/chapters", projectCreateChapterAction.Invoke)
	project.POST("/:id/chapters/:chapter_internal_id/reorder", projectReorderChapterAction.Invoke)
	project.POST("/:id/chapters/:chapter_internal_id", projectUpdateChapterAction.Invoke)
	project.DELETE("/:id/chapters/:chapter_internal_id", projectDeleteChapterAction.Invoke)

	project.POST("/:id/characters", projectCreateCharacterAction.Invoke)
	project.POST("/:id/characters/:character_internal_id", projectUpdateCharacterAction.Invoke)
	project.DELETE("/:id/characters/:character_internal_id", projectDeleteCharacterAction.Invoke)

	project.POST("/:id/locations", projectCreateLocationAction.Invoke)
	project.POST("/:id/locations/:location_internal_id", projectUpdateLocationAction.Invoke)
	project.DELETE("/:id/locations/:location_internal_id", projectDeleteLocationAction.Invoke)

	project.POST("/:id/metadata", projectUpdateMetadataAction.Invoke)

	project.POST("/:id/publish", projectPublishAction.Invoke)
	project.POST("/:id/republish", projectRepublishAction.Invoke)
}
