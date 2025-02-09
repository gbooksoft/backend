package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	middleware.RegisterLimit("GET", "/genres/", 4)
}

func genreRoutes(e *gin.Engine) {
	genre := e.Group("/genres")
	genre.Use(middleware.AuthRequired())

	genre.GET("", genreGetAllAction.Invoke)
}
