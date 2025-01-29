package action

import (
	"github.com/gbooksoft/backend/internal/api/action/response"
	"github.com/gbooksoft/backend/internal/domain/account"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Account struct {
	service *account.Service
}

func NewAccount(service *account.Service) *Account {
	return &Account{
		service: service,
	}
}

type confirmationEmailChangeRequest struct {
	Token string `json:"token"`
}

func (a *Account) ConfirmEmailChange(ctx *gin.Context) {
	var req confirmationEmailChangeRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	err := a.service.ConfirmEmailChange(req.Token)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
}
