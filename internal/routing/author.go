package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func authorRoutes(e *gin.Engine) {
	author := e.Group("/authors")

	author.Use(middleware.AuthRequired())

	author.GET("/:username", authorGetByUsernameAction.Invoke)
}
