package application

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/internal/domain/password"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gbooksoft/backend/pkg/token"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"os"
)

type SignUpCommand struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type SignUpHandlerResponse struct {
	UserId          string
	ValidationError *password.PasswordValidationResult
}

type SignUpHandler struct {
	db          *sql.DB
	emailFacade *email.Facade
}

func NewSignUpHandler(db *sql.DB, emailFacade *email.Facade) *SignUpHandler {
	return &SignUpHandler{
		db:          db,
		emailFacade: emailFacade,
	}
}

func (h *SignUpHandler) Handle(ctx context.Context, cmd *SignUpCommand) (*SignUpHandlerResponse, error) {
	userEmail := cmd.Email
	uPassword := cmd.Password
	firstName := cmd.FirstName
	lastName := cmd.LastName

	var exists sql.NullBool
	err := h.db.QueryRow("SELECT 1 FROM users WHERE email = ?", userEmail).Scan(&exists)

	if err == nil {
		return nil, gerr.NewBusinessLogicErr("can not register user", gerr.AuthorizationUserAlreadyExists)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		log.Error().Err(err).Msgf("can not register user %s", userEmail)
		return nil, err
	}

	validationResult := password.ValidatePassword(uPassword)

	if len(validationResult.Errors) > 0 {
		return &SignUpHandlerResponse{
			UserId:          "",
			ValidationError: validationResult,
		}, nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(uPassword), bcrypt.DefaultCost)

	if err != nil {
		log.Error().Err(err).Msgf("can hash password while register user %s", userEmail)
		return nil, err
	}

	confirmationHash, err := token.GenerateRandHash(64)

	if err != nil {
		log.Error().Err(err).Msgf("can not generate confirmation hash for user %s", userEmail)
		return nil, err
	}

	tx, err := h.db.Begin()
	if err != nil {
		log.Error().Err(err).Msgf("can not start transaction")
		return nil, err
	}
	defer tx.Rollback()

	userId := uuid.New().String()
	_, err = tx.Exec("INSERT INTO users (id, email, password, first_name, last_name, confirmation_hash ) VALUES (?, ?, ?, ?, ?, ?)", userId, userEmail, passwordHash, firstName, lastName, confirmationHash)

	if err != nil {
		log.Error().Err(err).Msgf("can not create user while register %s", userEmail)
		return nil, err
	}

	link, err := h.generateConfirmationLink(userId, confirmationHash)

	if err != nil {
		log.Error().Err(err).Msgf("can not generate confirmation link for user %s", userEmail)
		return nil, err
	}

	err = h.emailFacade.RegistrationConfirmationTx(tx, userEmail, email.RegistrationConfirmationData{
		FirstName:        firstName,
		LastName:         lastName,
		ConfirmationLink: link,
	})

	if err != nil {
		log.Error().Err(err).Msgf("can not add email to queue while register %s", userEmail)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msgf("can not commit transaction while register %s", userEmail)
		return nil, err
	}

	return &SignUpHandlerResponse{
		UserId:          userId,
		ValidationError: nil,
	}, nil
}

func (s *SignUpHandler) generateConfirmationLink(userId, hash string) (string, error) {
	u, err := url.Parse(os.Getenv("WEB_FRONT_REGISTRATION_CONFIRMATION_URL"))
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("user_id", userId)
	params.Add("token", hash)

	u.RawQuery = params.Encode()

	return u.String(), nil
}
