package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	middleware.RegisterLimit("GET", "/projects/", 10)
	middleware.RegisterLimit("POST", "/projects/", 10)

	middleware.RegisterLimit("GET", "/projects/recent", 10)

	middleware.RegisterLimit("GET", "/projects/:id/access", 10)
	middleware.RegisterLimit("GET", "/projects/:id/metadata", 10)

	middleware.RegisterLimit("GET", "/projects/:id/menu", 10)
	middleware.RegisterLimit("POST", "/projects/:id/menu", 10)
	middleware.RegisterLimit("POST", "/projects/:id/menu/:item_id/toggle-expand", 10)
	middleware.RegisterLimit("POST", "/projects/:id/menu/:item_id", 10)
	middleware.RegisterLimit("DELETE", "/projects/:id/menu/:item_id", 10)

	middleware.RegisterLimit("GET", "/projects/:id/editors", 10)
	middleware.RegisterLimit("POST", "/projects/:id/editors/:item_id", 10)

	middleware.RegisterLimit("GET", "/projects/:id/characters", 10)
	middleware.RegisterLimit("POST", "/projects/:id/characters/:item_id", 10)

	middleware.RegisterLimit("GET", "/projects/:id/locations", 10)
	middleware.RegisterLimit("POST", "/projects/:id/locations/:item_id", 10)

	middleware.RegisterLimit("GET", "/projects/:id", 10)
	middleware.RegisterLimit("POST", "/projects/:id", 10)
}

func projectRoutes(e *gin.Engine) {
	project := e.Group("/projects")

	project.Use(middleware.AuthRequired())

	project.POST("", projectCreateAction.Invoke)
	project.GET("", projectGetAllAction.Invoke)

	project.GET("/recent", projectGetRecentAction.Invoke)

	project.GET("/:id/access", projectCanAccessAction.Invoke)
	project.GET("/:id/metadata", projectGetMetadataAction.Invoke)
	project.POST("/:id/metadata", projectUpdateMetadataAction.Invoke)

	project.GET("/:id/menu", projectGetMenuAction.Invoke)
	project.POST("/:id/menu", projectAddMenuItemAction.Invoke)
	project.POST("/:id/menu/reorder", projectReorderMenuAction.Invoke)
	project.POST("/:id/menu/:item_id/toggle-expand", projectToggleExpandMenuItemAction.Invoke)
	project.POST("/:id/menu/:item_id", projectUpdateMenuItemAction.Invoke)
	project.DELETE("/:id/menu/:item_id", projectDeleteMenuItemAction.Invoke)

	project.GET("/:id/editors", projectGetEditorsAction.Invoke)
	project.POST("/:id/editors/:item_id", projectSaveEditorStateAction.Invoke)

	project.GET("/:id/characters", projectGetCharactersAction.Invoke)
	project.POST("/:id/characters/:item_id", projectUpdateCharacterAction.Invoke)

	project.GET("/:id/locations", projectGetLocationsAction.Invoke)
	project.POST("/:id/locations/:item_id", projectUpdateLocationAction.Invoke)

	project.GET("/:id", projectGetAction.Invoke)
	project.POST("/:id", projectUpdateAction.Invoke)
}
