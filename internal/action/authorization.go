package action

import (
	"github.com/gbooksoft/backend/internal/action/response"
	"github.com/gbooksoft/backend/internal/domain/authorization"
	"github.com/gin-gonic/gin"
)

type loggedResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresAt             int64  `json:"expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at"`
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
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (a *Authorize) SignIn(ctx *gin.Context) {
	var req signInRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	token, err := a.service.SignInUser(req.Email, req.Password)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	rsp := loggedResponse{
		AccessToken:           token.AccessToken,
		ExpiresAt:             token.AccessTokenExp,
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiresAt: token.RefreshTokenExp,
	}

	response.Json(ctx.Writer).StatusCode(200).Struct(rsp).Write()
	return
}

type signUpRequest struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (a *Authorize) SignUp(ctx *gin.Context) {
	var req signUpRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}
}
