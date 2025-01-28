package action

import (
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/internal/api/action/response"
	"github.com/gbooksoft/backend/internal/domain/account"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"path/filepath"
)

type Account struct {
	service *account.Service
}

func NewAccount(service *account.Service) *Account {
	return &Account{
		service: service,
	}
}

type updateSettingsRequest struct {
	Email     string `json:"email" validate:"omitempty,email" form:"email"`
	FirstName string `json:"first_name" validate:"omitempty" form:"first_name"`
	LastName  string `json:"last_name" validate:"omitempty" form:"last_name"`
}

type updateSettingsResponse struct {
	Email     bool `json:"email"`
	FirstName bool `json:"first_name"`
	LastName  bool `json:"last_name"`
	AvatarUrl bool `json:"avatar_url"`
}

func (a *Account) UpdateSettings(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	var req updateSettingsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err = validate.Struct(req)

	if err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	header, err := ctx.FormFile("avatar")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		response.Json(ctx.Writer).StatusCode(http.StatusBadRequest).Struct(struct {
			Error string `json:"error"`
		}{
			Error: "Failed to get avatar",
		}).Write()
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

	if header != nil {
		const maxFileSize = 2 << 20 // 2 MB
		if header.Size > maxFileSize {
			response.Json(ctx.Writer).StatusCode(http.StatusBadRequest).Struct(struct {
				Error string `json:"error"`
			}{
				Error: fmt.Sprintf("Avatar file is too large, max %vMB", maxFileSize>>20),
			}).Write()
			return
		}

		file, err := header.Open()
		if err != nil {
			response.Json(ctx.Writer).StatusCode(http.StatusBadRequest).Struct(struct {
				Error string `json:"error"`
			}{
				Error: "Failed to handle avatar",
			}).Write()
			return
		}
		defer file.Close()

		fileExt := filepath.Ext(header.Filename)
		s3FilePath := fmt.Sprintf("users/%s/avatar%s", claims.Subject, fileExt)

		// Загрузка файла на S3
		err = a.service.UploadAvatarToS3(file, s3FilePath)

		if err != nil {
			response.Json(ctx.Writer).StatusCode(http.StatusBadRequest).Struct(struct {
				Error string `json:"error"`
			}{
				Error: "Failed to upload avatar",
			}).Write()
			return
		}

		settings.AvatarPath = &s3FilePath
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
		AvatarUrl: updatedFields.AvatarPath,
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
