package action

import (
	"github.com/gbooksoft/backend/internal/api/action/response"
	"github.com/gbooksoft/backend/internal/domain/account/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type AccountConfirmEmailChange struct {
}

func NewAccountConfirmEmailChange() *AccountConfirmEmailChange {
	return new(AccountConfirmEmailChange)
}

type confirmationEmailChangeRequest struct {
	Token string `json:"token"`
}

func (a *AccountConfirmEmailChange) Invoke(ctx *gin.Context) {
	var req confirmationEmailChangeRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	cmd := &application.ConfirmEmailChangeCommand{
		Token: req.Token,
	}

	_, err := mediatr.Send[*application.ConfirmEmailChangeCommand, *application.ConfirmEmailChangeCommandHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
}
