package action

import (
	"github.com/gbooksoft/backend/internal/domain/authorization/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type AuthForgotPassword struct {
}

func NewAuthForgotPassword() *AuthForgotPassword {
	return new(AuthForgotPassword)
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

func (a *AuthForgotPassword) Invoke(ctx *gin.Context) {
	var req forgotPasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}
	cmd := &application.SendForgotPasswordNotificationCommand{
		Email: req.Email,
	}

	_, err := mediatr.Send[*application.SendForgotPasswordNotificationCommand, *application.SendForgotPasswordNotificationHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
