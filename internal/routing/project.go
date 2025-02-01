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
	project.POST("/:id/menu/:item_id/toggle-expand", projectToggleExpandMenuItemAction.Invoke)
	project.POST("/:id/menu/:item_id", actionProject.UpdateMenuItem)
	project.DELETE("/:id/menu/:item_id", actionProject.DeleteMenuItem)

	project.GET("/:id/editors", actionProject.GetEditors)
	project.POST("/:id/editors/:item_id", actionProject.SaveEditorState)

	project.GET("/:id/characters", actionProject.GetCharacters)
	project.POST("/:id/characters/:item_id", actionProject.UpdateCharacter)

	project.GET("/:id/locations", actionProject.GetLocations)
	project.POST("/:id/locations/:item_id", actionProject.UpdateLocation)

	project.GET("/:id", actionProject.GetProject)
	project.POST("/:id", actionProject.UpdateProject)
}
