package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	middleware.RegisterLimit("POST", "/auth/authorize", 10)
	middleware.RegisterLimit("POST", "/auth/register", 10)
	middleware.RegisterLimit("POST", "/auth/refresh", 10)
	middleware.RegisterLimit("POST", "/auth/logout", 10)
	middleware.RegisterLimit("POST", "/auth/confirm-registration", 10)
}

func authRoutes(e *gin.Engine) {
	oauth := e.Group("/auth")

	oauth.POST("/authorize", actionAuthorize.SignIn)
	oauth.POST("/register", actionAuthorize.SignUp)
	oauth.POST("/refresh", actionAuthorize.Refresh)
	oauth.POST("/logout", actionAuthorize.Logout)
	oauth.POST("/confirm-registration", actionAuthorize.ConfirmRegistration)
}
