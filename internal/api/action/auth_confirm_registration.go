package action

import (
	"github.com/gbooksoft/backend/internal/domain/authorization/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type AuthConfirmRegistration struct {
}

func NewAuthConfirmRegistration() *AuthConfirmRegistration {
	return new(AuthConfirmRegistration)
}

type confirmationRegistrationRequest struct {
	UserId string `json:"user_id"`
	Hash   string `json:"token"`
}

func (a *AuthConfirmRegistration) Invoke(ctx *gin.Context) {
	var req confirmationRegistrationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	cmd := &application.ConfirmRegistrationCommand{
		UserId: req.UserId,
		Hash:   req.Hash,
	}

	_, err := mediatr.Send[*application.ConfirmRegistrationCommand, *application.ConfirmRegistrationHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
}
