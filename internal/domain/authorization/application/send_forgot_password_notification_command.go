package application

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gbooksoft/backend/pkg/token"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
)

type SendForgotPasswordNotificationCommand struct {
	Email string
}

type SendForgotPasswordNotificationHandlerResponse struct{}

type SendForgotPasswordNotificationHandler struct {
	db          *sql.DB
	emailFacade *email.Facade
}

func NewSendForgotPasswordNotificationHandler(db *sql.DB, emailFacade *email.Facade) *SendForgotPasswordNotificationHandler {
	return &SendForgotPasswordNotificationHandler{
		db:          db,
		emailFacade: emailFacade,
	}
}

func (h *SendForgotPasswordNotificationHandler) Handle(ctx context.Context, cmd *SendForgotPasswordNotificationCommand) (*SendForgotPasswordNotificationHandlerResponse, error) {
	userEmail := cmd.Email

	tx, err := h.db.Begin()
	if err != nil {
		log.Error().Err(err).Msgf("can not start transaction while sending forgot password notification")
		return nil, err
	}
	defer tx.Rollback()

	var exists sql.NullBool
	err = tx.QueryRow("SELECT 1 FROM users WHERE email = ?", userEmail).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, gerr.NewBusinessLogicErr("user not found", gerr.AuthorizationUserNotFound)
		}

		log.Error().Err(err).Msgf("user %s not found", userEmail)
		return nil, err
	}

	resetPasswordHash, err := token.GenerateRandHash(64)

	if err != nil {
		log.Error().Err(err).Msgf("can not generate reset password hash for user %s", userEmail)
		return nil, err
	}

	_, err = tx.Exec("UPDATE users SET reset_password_hash = ?, reset_password_hash_created_at = NOW() WHERE email = ?", resetPasswordHash, userEmail)

	if err != nil {
		log.Error().Err(err).Msgf("can not save reset password hash for user %s", userEmail)
		return nil, err
	}

	link, err := h.generateResetPasswordLink(resetPasswordHash)

	if err != nil {
		log.Error().Err(err).Msgf("can not generate reset link")
	}

	err = h.emailFacade.ForgotPasswordTx(tx, userEmail, email.ForgotPasswordData{
		ResetPasswordLink: link,
	})

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msgf("can not commit transaction while sending forgot password notification %s", userEmail)
		return nil, err
	}

	return nil, nil
}

func (h *SendForgotPasswordNotificationHandler) generateResetPasswordLink(hash string) (string, error) {
	u, err := url.Parse(os.Getenv("WEB_FRONT_RESET_PASSWORD_URL"))
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("token", hash)

	u.RawQuery = params.Encode()

	return u.String(), nil
}
