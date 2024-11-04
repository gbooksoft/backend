package routing

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func authRoutes(e *gin.Engine) {
	oauth := e.Group("/auth")

	oauth.POST("/authorize", actionAuthorize.SignIn)
	oauth.POST("/register", actionAuthorize.SignUp)
}
