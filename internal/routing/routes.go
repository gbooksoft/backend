package routing

import (
	"github.com/gin-gonic/gin"
)

func Add(e *gin.Engine) {
	initializeDependencies()

	authRoutes(e)
	accountRoutes(e)
	projectRoutes(e)
}
