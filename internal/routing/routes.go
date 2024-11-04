package routing

import (
	"github.com/gin-gonic/gin"
)

func Add(e *gin.Engine) {
	initializeDependencies()

	infoRoutes(e)
	authRoutes(e)
}
