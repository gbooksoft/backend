package routing

import (
	"github.com/gbooksoft/backend/internal/action"
	"github.com/gin-gonic/gin"
)

func infoRoutes(e *gin.Engine) {
	e.GET("/about", action.About)
}
