package action

import (
	"github.com/gbooksoft/backend/internal/domain/authorization/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type AuthResetPassword struct {
}

func NewAuthResetPassword() *AuthResetPassword {
	return new(AuthResetPassword)
}

type resetPasswordRequest struct {
	Hash     string `json:"token"`
	Password string `json:"password"`
}
type resetPasswordErrorResponse struct {
	PasswordErrors []string `json:"password_errors"`
}

func (a *AuthResetPassword) Invoke(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	cmd := &application.ResetPasswordCommand{
		Hash:     req.Hash,
		Password: req.Password,
	}

	resp, err := mediatr.Send[*application.ResetPasswordCommand, *application.ResetPasswordHandlerResponse](ctx, cmd)

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

	ctx.JSON(http.StatusOK, 1)
	return
}
