package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func Add(e *gin.Engine) {
	e.GET("/logo_36_transparent", func(c *gin.Context) {
		c.File("static/images/logo_36_transparent.png")
	})
	e.Use(middleware.RateLimiterMiddleware())

	initializeDependencies()

	authRoutes(e)
	accountRoutes(e)
	projectRoutes(e)
	genreRoutes(e)
	bookRoutes(e)
	authorRoutes(e)
}
