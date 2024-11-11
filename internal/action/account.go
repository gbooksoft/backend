package action

import (
	"errors"
	"github.com/gbooksoft/backend/internal/action/response"
	"github.com/gbooksoft/backend/internal/domain/account"
	"github.com/gbooksoft/backend/internal/domain/password"
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

type updateSettingsRequest struct {
	Email     string `json:"email" validate:"omitempty,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type updateSettingsResponse struct {
	Email     bool `json:"email"`
	FirstName bool `json:"first_name"`
	LastName  bool `json:"last_name"`
}

func (a *Account) UpdateSettings(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	var req updateSettingsRequest
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

	updatedFields, err := a.service.Update(claims.Subject, settings)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Struct(updateSettingsResponse{
		Email:     updatedFields.Email,
		FirstName: updatedFields.FirstName,
		LastName:  updatedFields.LastName,
	}).Write()
	return
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

type updatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}

type updatePasswordPasswordValidationErrorResponse struct {
	PasswordErrors []string `json:"password_errors"`
}

func (a *Account) UpdatePassword(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	var req updatePasswordRequest
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

	err = a.service.UpdatePassword(claims.Subject, req.CurrentPassword, req.NewPassword)

	if err != nil {
		var passErr *password.PasswordValidationResult
		if errors.As(err, &passErr) {
			response.Json(ctx.Writer).StatusCode(http.StatusBadRequest).Struct(updatePasswordPasswordValidationErrorResponse{
				PasswordErrors: passErr.Errors,
			}).Write()
			return
		}

		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
	return
}
