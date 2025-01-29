package routing

import (
	"github.com/gbooksoft/backend/internal/routing/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	middleware.RegisterLimit("GET", "/account/profile", 10)
	middleware.RegisterLimit("POST", "/account/settings", 10)
	middleware.RegisterLimit("POST", "/account/settings/password", 10)
	middleware.RegisterLimit("POST", "/account/confirm-email-change", 10)
}

func accountRoutes(e *gin.Engine) {
	account := e.Group("/account")
	accountWithoutAuth := e.Group("/account")

	account.Use(middleware.AuthRequired())

	account.GET("/profile", accountGetProfileAction.Invoke)
	account.POST("/settings", accountUpdateSettingsAction.Invoke)
	account.POST("/settings/password", accountUpdatePasswordAction.Invoke)
	accountWithoutAuth.POST("/confirm-email-change", accountConfirmEmailChangeAction.Invoke)
}
