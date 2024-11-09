package action

import (
	"github.com/gbooksoft/backend/internal/action/response"
	"github.com/gbooksoft/backend/internal/domain/account"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func (a *Account) GetProfile(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	email, firstName, lastName, err := a.service.GetUser(claims.Subject)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	ctx.JSON(200, gin.H{
		"user_id":    claims.Subject,
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
	})
}

type updateSettings struct {
	Email     string `json:"email" validate:"omitempty,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

func (a *Account) UpdateSettings(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	var req updateSettings
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err = validate.Struct(req)

	if err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	settings := account.UpdateFields{}

	if req.Email != "" {
		settings.Email = &req.Email
	}

	if req.FirstName != "" {
		settings.FirstName = &req.FirstName
	}

	if req.LastName != "" {
		settings.LastName = &req.LastName
	}

	if req.Password != "" {
		settings.Password = &req.Password
	}

	err = a.service.Update(claims.Subject, settings)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
	return
}
