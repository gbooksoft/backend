package account

import (
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/rs/zerolog/log"
	"time"
)

type Service struct {
	db          *sql.DB
	emailFacade *email.Facade
}

func NewService(db *sql.DB, emailFacade *email.Facade) *Service {
	return &Service{
		db:          db,
		emailFacade: emailFacade,
	}
}

func (s *Service) ConfirmEmailChange(token string) error {
	var createdAtStr string
	var newEmail string
	var userId string

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRow("SELECT created_at, new_email, user_id FROM email_change_requests WHERE token = ?", token).Scan(&createdAtStr, &newEmail, &userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gerr.NewBusinessLogicErr("token not found", gerr.AccountEmailChangeRequestNotFound)
		}

		log.Error().Err(err).Msgf("failed to find email change with token %s", token)
		return err
	}

	layout := "2006-01-02 15:04:05"
	createdAt, err := time.Parse(layout, createdAtStr)
	if err != nil {
		return err
	}

	if time.Since(createdAt) >= 15*time.Minute {
		return gerr.NewBusinessLogicErr("token expired", gerr.AccountEmailChangeRequestExpired)
	}

	_, err = tx.Exec("UPDATE users SET email = ? WHERE id = ?", newEmail, userId)
	if err != nil {
		log.Error().Err(err).Msg("failed to confirm user email")
		return err
	}

	_, err = tx.Exec("DELETE FROM email_change_requests WHERE token = ?", token)

	if err != nil {
		log.Error().Err(err).Msg("failed to delete email change request")
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *Service) existsUserWithEmail(email string) bool {
	var exists sql.NullBool
	err := s.db.QueryRow("SELECT 1 FROM users WHERE email = ?", email).Scan(&exists)

	if err == nil {
		return true
	}

	return false
}
