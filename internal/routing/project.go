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
	project.GET("", actionProject.GetAllProjects)

	project.GET("/recent", actionProject.GetRecent)

	project.GET("/:id/access", actionProject.CanAccess)
	project.GET("/:id/metadata", actionProject.GetMetadata)
	project.POST("/:id/metadata", actionProject.CreateOrUpdateMetadata)

	project.GET("/:id/menu", actionProject.GetProjectMenu)
	project.POST("/:id/menu", actionProject.AddMenuItem)
	project.POST("/:id/menu/:item_id/toggle-expand", actionProject.ToggleExpandMenuItem)
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
