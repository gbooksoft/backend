package application

import (
	"context"
	"database/sql"
	"github.com/gbooksoft/backend/pkg/filesystem"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"time"
)

type UploadAvatarCommand struct {
	File       multipart.File
	S3FilePath string
}

type UploadAvatarHandlerResponse struct {
}

type UploadAvatarHandler struct {
	db *sql.DB
}

func NewUploadAvatarHandler(db *sql.DB) *UploadAvatarHandler {
	return &UploadAvatarHandler{
		db: db,
	}
}

func (h *UploadAvatarHandler) Handle(ctx context.Context, cmd *UploadAvatarCommand) (*UploadAvatarHandlerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := filesystem.UploadMultipartFile(ctx, cmd.S3FilePath, cmd.File)

	if err != nil {
		log.Error().Err(err).Msg("can not upload file")
		return nil, err
	}

	return nil, nil
}
