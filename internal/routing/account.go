package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func accountRoutes(e *gin.Engine) {
	account := e.Group("/account")
	accountWithoutAuth := e.Group("/account")

	account.Use(middleware.AuthRequired())

	account.GET("/profile", actionAccount.GetProfile)
	account.POST("/settings", actionAccount.UpdateSettings)
	accountWithoutAuth.POST("/confirm-email-change", actionAccount.ConfirmEmailChange)
}
