package action

import (
	"github.com/gin-gonic/gin"
	"os"
)

const refreshTokenCookieName = "refreshToken"

func setRefreshTokenCookie(ctx *gin.Context, refreshToken string) {
	ctx.SetCookie(refreshTokenCookieName, refreshToken, 3600*24*7, "/", os.Getenv("WEB_FRONT_DOMAIN"), false, true)
}

func deleteRefreshTokenCookie(ctx *gin.Context) {
	ctx.SetCookie(refreshTokenCookieName, "", -1, "/", os.Getenv("WEB_FRONT_DOMAIN"), false, true)
}

func getRefreshTokenFromCookie(ctx *gin.Context) (string, error) {
	rt, err := ctx.Cookie(refreshTokenCookieName)

	if err != nil {
		return "", err
	}

	return rt, nil
}
