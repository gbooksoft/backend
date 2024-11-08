package action

import (
	"github.com/gbooksoft/backend/internal/action/response"
	"github.com/gbooksoft/backend/internal/domain/authorization"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"os"
)

const refreshTokenCookieName = "refreshToken"

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

type signInRequest struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (a *Authorize) SignIn(ctx *gin.Context) {
	var req signInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(req)

	if err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	token, err := a.service.SignInUser(req.Email, req.Password)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	rsp := atResponse{
		AccessToken: token.AccessToken,
	}

	a.setRefreshTokenCookie(ctx, token.RefreshToken)

	response.Json(ctx.Writer).StatusCode(200).Struct(rsp).Write()
	return
}

type signUpRequest struct {
	Email     string `json:"email" binding:"required" validate:"email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
}
type signUpResponse struct {
	UserId string `json:"user_id"`
}

func (a *Authorize) SignUp(ctx *gin.Context) {
	var req signUpRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(req)

	if err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	userId, err := a.service.Register(req.Email, req.Password, req.FirstName, req.LastName)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(201).Struct(signUpResponse{UserId: userId}).Write()
	return
}

func (a *Authorize) Refresh(ctx *gin.Context) {
	rt, err := a.getRefreshTokenFromCookie(ctx)

	if err != nil {
		response.Json(ctx.Writer).StatusCode(http.StatusUnauthorized).ErrMsg("refresh token not found in cookie").Write()
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
		response.Json(ctx.Writer).StatusCode(http.StatusUnauthorized).ErrMsg("refresh token not found in cookie").Write()
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
