package action

import (
	"github.com/gbooksoft/backend/internal/domain/authorization/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type AuthSignIn struct {
}

func NewAuthSignIn() *AuthSignIn {
	return new(AuthSignIn)
}

type signInRequest struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (a *AuthSignIn) Invoke(ctx *gin.Context) {
	var req signInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(req)

	if err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	cmd := &application.SignInCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := mediatr.Send[*application.SignInCommand, *application.SignInHandlerResponse](ctx, cmd)

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
