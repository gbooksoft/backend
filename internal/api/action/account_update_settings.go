package action

import (
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/internal/api/action/response"
	"github.com/gbooksoft/backend/internal/domain/account/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
	"path/filepath"
)

type AccountUpdateSettings struct {
}

func NewAccountUpdateSettings() *AccountUpdateSettings {
	return new(AccountUpdateSettings)
}

type updateSettingsRequest struct {
	Email     string `json:"email" validate:"omitempty,email" form:"email"`
	FirstName string `json:"first_name" validate:"omitempty" form:"first_name"`
	LastName  string `json:"last_name" validate:"omitempty" form:"last_name"`
}

func (a *AccountUpdateSettings) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(gerr.NewValidationError(err))
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
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	header, err := ctx.FormFile("avatar")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get avatar",
		})
		return
	}

	settings := application.UpdateFields{}

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
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Avatar file is too large, max %vMB", maxFileSize>>20),
			})
			return
		}

		file, err := header.Open()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to handle avatar",
			})
			return
		}
		defer file.Close()

		fileExt := filepath.Ext(header.Filename)
		s3FilePath := fmt.Sprintf("users/%s/avatar%s", claims.Subject, fileExt)

		cmd := &application.UploadAvatarCommand{
			File:       file,
			S3FilePath: s3FilePath,
		}

		_, err = mediatr.Send[*application.UploadAvatarCommand, *application.UploadAvatarHandlerResponse](ctx, cmd)

		if err != nil {
			ctx.Error(err)
			return
		}

		settings.AvatarPath = &s3FilePath
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
		"avatar_url": updatedFields.AvatarPath,
	})
	return
}
