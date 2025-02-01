package action

import (
	"github.com/gbooksoft/backend/internal/domain/authorization/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type AuthLogout struct {
}

func NewAuthLogout() *AuthLogout {
	return new(AuthLogout)
}

func (a *AuthLogout) Invoke(ctx *gin.Context) {
	rt, err := getRefreshTokenFromCookie(ctx)

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error_message": "refresh token not found in cookie",
		})
		return
	}
	cmd := &application.LogoutCommand{
		RefreshToken: rt,
	}

	_, err = mediatr.Send[*application.LogoutCommand, *application.LogoutHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	deleteRefreshTokenCookie(ctx)
	ctx.JSON(http.StatusOK, 1)
	return
}
