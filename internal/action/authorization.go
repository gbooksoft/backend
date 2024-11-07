package action

import (
	"github.com/gbooksoft/backend/internal/action/response"
	"github.com/gbooksoft/backend/internal/domain/authorization"
	"github.com/gin-gonic/gin"
	"os"
)

type loggedResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *Authorize) SignIn(ctx *gin.Context) {
	var req signInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	token, err := a.service.SignInUser(req.Email, req.Password)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	rsp := loggedResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	ctx.SetCookie("refreshToken", token.RefreshToken, 3600*24*7, "/", os.Getenv("APP_DOMAIN"), true, true)

	response.Json(ctx.Writer).StatusCode(200).Struct(rsp).Write()
	return
}

type signUpRequest struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
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

	userId, err := a.service.Register(req.Email, req.Password)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(201).Struct(signUpResponse{UserId: userId}).Write()
	return
}
