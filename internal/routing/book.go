package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func bookRoutes(e *gin.Engine) {
	book := e.Group("/books")

	book.Use(middleware.AuthRequired())
}
