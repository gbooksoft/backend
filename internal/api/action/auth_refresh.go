package action

import (
	"github.com/gbooksoft/backend/internal/domain/authorization/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type AuthRefresh struct {
}

func NewAuthRefresh() *AuthRefresh {
	return new(AuthRefresh)
}

func (a *AuthRefresh) Invoke(ctx *gin.Context) {
	rt, err := getRefreshTokenFromCookie(ctx)

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error_message": "refresh token not found"})
		return
	}

	cmd := &application.RefreshCommand{
		RefreshToken: rt,
	}

	resp, err := mediatr.Send[*application.RefreshCommand, *application.RefreshHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	setRefreshTokenCookie(ctx, resp.RefreshToken)

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": resp.AccessToken,
	})
	return
}
