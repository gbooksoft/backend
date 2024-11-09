package action

import (
	"github.com/gbooksoft/backend/internal/action/response"
	"github.com/gbooksoft/backend/internal/domain/account"
	"github.com/gin-gonic/gin"
)

type Account struct {
	service *account.Service
}

func NewAccount(service *account.Service) *Account {
	return &Account{
		service: service,
	}
}

func (a *Account) GetProfile(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	email, err := a.service.GetUserEmail(claims.Subject)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	ctx.JSON(200, gin.H{
		"user_id":    claims.Subject,
		"email":      email,
		"first_name": claims.FirstName,
		"last_name":  claims.LastName,
	})
}
