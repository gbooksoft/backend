package routing

import (
	"github.com/gin-gonic/gin"
)

func Add(e *gin.Engine) {
	infoRoutes(e)
	oAuthRoutes(e)

}
