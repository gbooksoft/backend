package application

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/rs/zerolog/log"
)

type ConfirmRegistrationCommand struct {
	UserId string
	Hash   string
}

type ConfirmRegistrationHandlerResponse struct{}

type ConfirmRegistrationHandler struct {
	db *sql.DB
}

func NewConfirmRegistrationHandler(db *sql.DB) *ConfirmRegistrationHandler {
	return &ConfirmRegistrationHandler{db}
}

func (h *ConfirmRegistrationHandler) Handle(ctx context.Context, cmd *ConfirmRegistrationCommand) (*ConfirmRegistrationHandlerResponse, error) {
	userId := cmd.UserId
	hash := cmd.Hash

	tx, err := h.db.Begin()
	if err != nil {
		log.Error().Err(err).Msgf("can not start transaction while confirm registration")
		return nil, err
	}
	defer tx.Rollback()

	var exists sql.NullBool
	err = tx.QueryRow("SELECT 1 FROM users WHERE id = ? AND confirmation_hash = ? AND confirmed = 0", userId, hash).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, gerr.NewBusinessLogicErr("confirmation not found", gerr.AuthorizationConfirmationNoFound)
		}

		log.Error().Err(err).Msgf("confirmation not found for user %s and hash %s", userId, hash)
		return nil, err
	}

	_, err = tx.Exec("UPDATE users SET confirmation_hash = NULL, confirmed = 1 WHERE id = ?", userId)

	if err != nil {
		log.Error().Err(err).Msgf("can not confirm registration")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msgf("can not commit transaction while confirming register %s", userId)
		return nil, err
	}

	return nil, nil
}
