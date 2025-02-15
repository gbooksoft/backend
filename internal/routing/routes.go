package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func Add(e *gin.Engine) {
	e.Use(middleware.RateLimiterMiddleware())

	initializeDependencies()

	authRoutes(e)
	accountRoutes(e)
	projectRoutes(e)
	genreRoutes(e)
}
