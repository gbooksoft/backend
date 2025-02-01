package application

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/internal/domain/password"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type ResetPasswordCommand struct {
	Hash     string
	Password string
}

type ResetPasswordHandlerResponse struct {
	ValidationError *password.PasswordValidationResult
}

type ResetPasswordHandler struct {
	db          *sql.DB
	emailFacade *email.Facade
}

func NewResetPasswordHandler(db *sql.DB, emailFacade *email.Facade) *ResetPasswordHandler {
	return &ResetPasswordHandler{db: db, emailFacade: emailFacade}
}

func (h *ResetPasswordHandler) Handle(ctx context.Context, cmd *ResetPasswordCommand) (*ResetPasswordHandlerResponse, error) {
	newPassword := cmd.Password
	hash := cmd.Hash

	tx, err := h.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	validationResult := password.ValidatePassword(newPassword)

	if len(validationResult.Errors) > 0 {
		return &ResetPasswordHandlerResponse{ValidationError: validationResult}, nil
	}

	var hashCreatedAt sql.NullString
	var userId string
	var uEmail string
	var firstName string
	var lastName string
	query := `
SELECT id, email, first_name, last_name, reset_password_hash_created_at FROM users WHERE reset_password_hash  = ?
`
	err = tx.QueryRow(query, hash).Scan(&userId, &uEmail, &firstName, &lastName, &hashCreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, gerr.NewBusinessLogicErr("reset password not found", gerr.AuthorizationResetPasswordRequestNotFound)
		}

		return nil, err
	}

	if !hashCreatedAt.Valid {
		return nil, gerr.NewBusinessLogicErr("reset password not found", gerr.AuthorizationResetPasswordRequestNotFound)
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", hashCreatedAt.String)

	if err != nil {
		return nil, err
	}

	if time.Since(parsedTime) > 24*time.Hour {
		return nil, gerr.NewBusinessLogicErr("reset password not found", gerr.AuthorizationResetPasswordRequestNotFound)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		log.Error().Err(err).Msgf("can hash password while changing it")
		return nil, err
	}

	_, err = tx.Exec("UPDATE users SET password = ?, reset_password_hash_created_at = NULL, reset_password_hash = NULL, confirmed = 1, confirmation_hash = NULL WHERE id = ?", passwordHash, userId)

	if err != nil {
		log.Error().Err(err).Msgf("can update password while sql update")
		return nil, err
	}

	err = h.emailFacade.PasswordChangedTx(tx, uEmail, email.PasswordChangedData{
		FirstName: firstName,
		LastName:  lastName,
	})
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
