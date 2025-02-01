package action

import (
	"github.com/gbooksoft/backend/internal/domain/account/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type AccountUpdatePassword struct {
}

func NewAccountUpdatePassword() *AccountUpdatePassword {
	return new(AccountUpdatePassword)
}

type updatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}

func (a *AccountUpdatePassword) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(gerr.NewBusinessLogicErr("Unauthorized", gerr.AuthorizationUserNotFound))
		return
	}

	var req updatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err = validate.Struct(req)

	if err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	cmd := &application.UpdatePasswordCommand{
		UserId:          claims.Subject,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}

	resp, err := mediatr.Send[*application.UpdatePasswordCommand, *application.UpdatePasswordHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	if resp != nil && resp.ValidationError != nil {
		ctx.JSON(http.StatusBadRequest, struct {
			PasswordErrors []string `json:"password_errors"`
		}{
			PasswordErrors: resp.ValidationError.Errors,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
