package action

import (
	"github.com/gbooksoft/backend/internal/domain/account/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type AccountUpdateSettings struct {
}

func NewAccountUpdateSettings() *AccountUpdateSettings {
	return new(AccountUpdateSettings)
}

type updateSettingsRequest struct {
	Email     string `json:"email" validate:"omitempty,email" form:"email"`
	Username  string `json:"username" validate:"omitempty" form:"username"`
	FirstName string `json:"first_name" validate:"omitempty" form:"first_name"`
	LastName  string `json:"last_name" validate:"omitempty" form:"last_name"`
}

func (a *AccountUpdateSettings) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	var req updateSettingsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err = validate.Struct(req)

	if err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	settings := application.UpdateFields{}

	if req.Email != "" {
		settings.Email = &req.Email
	}

	if req.Username != "" {
		settings.Username = &req.Username
	}

	if req.FirstName != "" {
		settings.FirstName = &req.FirstName
	}

	if req.LastName != "" {
		settings.LastName = &req.LastName
	}

	cmd := &application.UpdateSettingsCommand{
		UserId:       claims.Subject,
		UpdateFields: settings,
	}

	resp, err := mediatr.Send[*application.UpdateSettingsCommand, *application.UpdateSettingsHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	updatedFields := resp.UpdatedFields

	ctx.JSON(http.StatusOK, gin.H{
		"email":      updatedFields.Email,
		"first_name": updatedFields.FirstName,
		"last_name":  updatedFields.LastName,
	})
	return
}
