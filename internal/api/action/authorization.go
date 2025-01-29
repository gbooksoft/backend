package action

import (
	"errors"
	"github.com/gbooksoft/backend/internal/api/action/response"
	"github.com/gbooksoft/backend/internal/domain/authorization"
	"github.com/gbooksoft/backend/internal/domain/password"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type atResponse struct {
	AccessToken string `json:"access_token"`
}

type Authorize struct {
	service *authorization.Service
}

func NewAuthorize(service *authorization.Service) *Authorize {
	return &Authorize{
		service: service,
	}
}

func (a *Authorize) Refresh(ctx *gin.Context) {
	rt, err := a.getRefreshTokenFromCookie(ctx)

	if err != nil {
		response.Json(ctx.Writer).StatusCode(http.StatusForbidden).ErrMsg("refresh token not found in cookie").Write()
		return
	}

	token, err := a.service.RefreshToken(rt)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	rsp := atResponse{
		AccessToken: token.AccessToken,
	}

	a.setRefreshTokenCookie(ctx, token.RefreshToken)

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Struct(rsp).Write()
	return
}

func (a *Authorize) Logout(ctx *gin.Context) {
	rt, err := a.getRefreshTokenFromCookie(ctx)

	if err != nil {
		response.Json(ctx.Writer).StatusCode(http.StatusForbidden).ErrMsg("refresh token not found in cookie").Write()
		return
	}

	err = a.service.Logout(rt)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	a.deleteRefreshTokenCookie(ctx)
	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
}

type confirmationRegistrationRequest struct {
	UserId string `json:"user_id"`
	Hash   string `json:"token"`
}

func (a *Authorize) ConfirmRegistration(ctx *gin.Context) {
	var req confirmationRegistrationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	err := a.service.ConfirmRegistration(req.UserId, req.Hash)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
}

func (a *Authorize) setRefreshTokenCookie(ctx *gin.Context, refreshToken string) {
	ctx.SetCookie(refreshTokenCookieName, refreshToken, 3600*24*7, "/", os.Getenv("WEB_FRONT_DOMAIN"), false, true)
}

func (a *Authorize) deleteRefreshTokenCookie(ctx *gin.Context) {
	ctx.SetCookie(refreshTokenCookieName, "refreshToken", -1, "/", os.Getenv("WEB_FRONT_DOMAIN"), false, true)
}

func (a *Authorize) getRefreshTokenFromCookie(ctx *gin.Context) (string, error) {
	rt, err := ctx.Cookie(refreshTokenCookieName)

	if err != nil {
		return "", err
	}

	return rt, nil
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

func (a *Authorize) SendForgotPasswordNotification(ctx *gin.Context) {
	var req forgotPasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	err := a.service.SendForgotPasswordNotification(req.Email)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
}

type resetPasswordRequest struct {
	Hash     string `json:"token"`
	Password string `json:"password"`
}
type resetPasswordErrorResponse struct {
	PasswordErrors []string `json:"password_errors"`
}

func (a *Authorize) ResetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	err := a.service.ResetPassword(req.Hash, req.Password)

	if err != nil {
		var passErr *password.PasswordValidationResult
		if errors.As(err, &passErr) {
			response.Json(ctx.Writer).StatusCode(http.StatusBadRequest).Struct(resetPasswordErrorResponse{
				PasswordErrors: passErr.Errors,
			}).Write()
			return
		}

		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
}
