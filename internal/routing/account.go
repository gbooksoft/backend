package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func accountRoutes(e *gin.Engine) {
	account := e.Group("/account")

	account.Use(middleware.AuthRequired())

	account.GET("/profile", actionAccount.GetProfile)
}
