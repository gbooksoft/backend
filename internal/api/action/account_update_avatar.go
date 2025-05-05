package action

import (
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/internal/domain/account/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/rs/zerolog/log"
	"net/http"
	"path/filepath"
)

type AccountUpdateAvatar struct {
}

func NewAccountUpdateAvatar() *AccountUpdateAvatar {
	return new(AccountUpdateAvatar)
}

func (a *AccountUpdateAvatar) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(gerr.NewBusinessLogicErr("Unauthorized", gerr.AuthorizationUserNotFound))
		return
	}

	header, err := ctx.FormFile("avatar")
	if (err != nil && !errors.Is(err, http.ErrMissingFile)) || header == nil {
		log.Err(err).Msg("avatar file is missing")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get avatar",
		})
		return
	}

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
		UserId:     claims.Subject,
	}

	resp, err := mediatr.Send[*application.UploadAvatarCommand, *application.UploadAvatarHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"avatar_url": resp.AvatarUrl,
	})
}
