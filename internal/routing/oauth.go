package routing

import (
	"github.com/gbooksoft/backend/internal/domain/authorization"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func oAuthRoutes(e *gin.Engine) {
	oAuthServer := authorization.CreateOAuthServer()

	oauth := e.Group("/oauth")
	oauth.POST("/token", func(ctx *gin.Context) {
		err := oAuthServer.HandleTokenRequest(ctx.Writer, ctx.Request)

		if err != nil {
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			ctx.Writer.Write([]byte(err.Error()))
		}
	})
}
