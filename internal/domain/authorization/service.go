package authorization

import (
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/internal/domain/password"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gbooksoft/backend/pkg/token"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"os"
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

type Token struct {
	AccessToken     string
	AccessTokenExp  int64
	RefreshToken    string
	RefreshTokenExp int64
}

func (s *Service) ConfirmRegistration(userId string, hash string) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Error().Err(err).Msgf("can not start transaction while confirm registration")
		return err
	}
	defer tx.Rollback()

	var exists sql.NullBool
	err = tx.QueryRow("SELECT 1 FROM users WHERE id = ? AND confirmation_hash = ? AND confirmed = 0", userId, hash).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gerr.NewBusinessLogicErr("confirmation not found", gerr.AuthorizationConfirmationNoFound)
		}

		log.Error().Err(err).Msgf("confirmation not found for user %s and hash %s", userId, hash)
		return err
	}

	_, err = tx.Exec("UPDATE users SET confirmation_hash = NULL, confirmed = 1 WHERE id = ?", userId)

	if err != nil {
		log.Error().Err(err).Msgf("can not confirm registration")
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msgf("can not commit transaction while confirming register %s", userId)
		return err
	}

	return nil
}

func (s *Service) SendForgotPasswordNotification(userEmail string) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Error().Err(err).Msgf("can not start transaction while sending forgot password notification")
		return err
	}
	defer tx.Rollback()

	var exists sql.NullBool
	err = tx.QueryRow("SELECT 1 FROM users WHERE email = ?", userEmail).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gerr.NewBusinessLogicErr("user not found", gerr.AuthorizationUserNotFound)
		}

		log.Error().Err(err).Msgf("user %s not found", userEmail)
		return err
	}

	resetPasswordHash, err := token.GenerateRandHash(64)

	if err != nil {
		log.Error().Err(err).Msgf("can not generate reset password hash for user %s", userEmail)
		return err
	}

	_, err = tx.Exec("UPDATE users SET reset_password_hash = ?, reset_password_hash_created_at = NOW() WHERE email = ?", resetPasswordHash, userEmail)

	if err != nil {
		log.Error().Err(err).Msgf("can not save reset password hash for user %s", userEmail)
		return err
	}

	link, err := s.generateResetPasswordLink(resetPasswordHash)

	if err != nil {
		log.Error().Err(err).Msgf("can not generate reset link")
	}

	err = s.emailFacade.ForgotPasswordTx(tx, userEmail, email.ForgotPasswordData{
		ResetPasswordLink: link,
	})

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msgf("can not commit transaction while sending forgot password notification %s", userEmail)
		return err
	}

	return nil
}

func (s *Service) generateResetPasswordLink(hash string) (string, error) {
	u, err := url.Parse(os.Getenv("WEB_FRONT_RESET_PASSWORD_URL"))
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("token", hash)

	u.RawQuery = params.Encode()

	return u.String(), nil
}

func (s *Service) ResetPassword(hash string, newPassword string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	validationResult := password.ValidatePassword(newPassword)

	if len(validationResult.Errors) > 0 {
		return validationResult
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
			return gerr.NewBusinessLogicErr("reset password not found", gerr.AuthorizationResetPasswordRequestNotFound)
		}

		return err
	}

	if !hashCreatedAt.Valid {
		return gerr.NewBusinessLogicErr("reset password not found", gerr.AuthorizationResetPasswordRequestNotFound)
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", hashCreatedAt.String)

	if err != nil {
		return err
	}

	if time.Since(parsedTime) > 24*time.Hour {
		return gerr.NewBusinessLogicErr("reset password not found", gerr.AuthorizationResetPasswordRequestNotFound)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		log.Error().Err(err).Msgf("can hash password while changing it")
		return err
	}

	_, err = tx.Exec("UPDATE users SET password = ?, reset_password_hash_created_at = NULL, reset_password_hash = NULL, confirmed = 1, confirmation_hash = NULL WHERE id = ?", passwordHash, userId)

	if err != nil {
		log.Error().Err(err).Msgf("can update password while sql update")
		return err
	}

	err = s.emailFacade.PasswordChangedTx(tx, uEmail, email.PasswordChangedData{
		FirstName: firstName,
		LastName:  lastName,
	})
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
