package action

import (
	"github.com/gbooksoft/backend/internal/domain/authorization/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type AuthSignUp struct {
}

func NewAuthSignUp() *AuthSignUp {
	return new(AuthSignUp)
}

type signUpRequest struct {
	Email     string `json:"email" binding:"required" validate:"email"`
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (a *AuthSignUp) Invoke(ctx *gin.Context) {
	var req signUpRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(req)

	if err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	cmd := &application.SignUpCommand{
		Email:     req.Email,
		Password:  req.Password,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	resp, err := mediatr.Send[*application.SignUpCommand, *application.SignUpHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	if resp.ValidationError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"password_errors": resp.ValidationError.Errors,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user_id": resp.UserId,
	})
	return
}
