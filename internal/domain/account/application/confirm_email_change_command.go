package application

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/rs/zerolog/log"
	"time"
)

type ConfirmEmailChangeCommand struct {
	Token string
}

type ConfirmEmailChangeCommandHandlerResponse struct {
}

type ConfirmEmailChangeCommandHandler struct {
	db *sql.DB
}

func NewConfirmEmailChangeCommandHandler(db *sql.DB) *ConfirmEmailChangeCommandHandler {
	return &ConfirmEmailChangeCommandHandler{
		db: db,
	}
}

func (h *ConfirmEmailChangeCommandHandler) Handle(ctx context.Context, cmd *ConfirmEmailChangeCommand) (*ConfirmEmailChangeCommandHandlerResponse, error) {
	var createdAtStr string
	var newEmail string
	var userId string

	token := cmd.Token

	tx, err := h.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRow("SELECT created_at, new_email, user_id FROM email_change_requests WHERE token = ?", token).Scan(&createdAtStr, &newEmail, &userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, gerr.NewBusinessLogicErr("token not found", gerr.AccountEmailChangeRequestNotFound)
		}

		log.Error().Err(err).Msgf("failed to find email change with token %s", token)
		return nil, err
	}

	layout := "2006-01-02 15:04:05"
	createdAt, err := time.Parse(layout, createdAtStr)
	if err != nil {
		return nil, err
	}

	if time.Since(createdAt) >= 15*time.Minute {
		return nil, gerr.NewBusinessLogicErr("token expired", gerr.AccountEmailChangeRequestExpired)
	}

	_, err = tx.Exec("UPDATE users SET email = ? WHERE id = ?", newEmail, userId)
	if err != nil {
		log.Error().Err(err).Msg("failed to confirm user email")
		return nil, err
	}

	_, err = tx.Exec("DELETE FROM email_change_requests WHERE token = ?", token)

	if err != nil {
		log.Error().Err(err).Msg("failed to delete email change request")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
