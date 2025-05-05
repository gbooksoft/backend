package application

import (
	"context"
	"database/sql"
	"github.com/gbooksoft/backend/pkg/filesystem"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"time"
)

type UploadCoverCommand struct {
	File       multipart.File
	S3FilePath string
	UserId     string
}

type UploadCoverHandlerResponse struct {
	CoverUrl string
}

type UploadCoverHandler struct {
	db *sql.DB
}

func NewUploadCoverHandler(db *sql.DB) *UploadCoverHandler {
	return &UploadCoverHandler{
		db: db,
	}
}

func (h *UploadCoverHandler) Handle(ctx context.Context, cmd *UploadCoverCommand) (*UploadCoverHandlerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := filesystem.UploadMultipartFile(ctx, cmd.S3FilePath, cmd.File)

	if err != nil {
		log.Error().Err(err).Msg("can not upload file")
		return nil, err
	}

	tx, err := h.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE users SET cover_path = ? WHERE id = ?", cmd.S3FilePath, cmd.UserId)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to update user's cover")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	coverUrl, err := filesystem.GetPresignedURL(cmd.S3FilePath)

	if err != nil {
		return nil, err
	}

	return &UploadCoverHandlerResponse{CoverUrl: coverUrl}, nil
}
