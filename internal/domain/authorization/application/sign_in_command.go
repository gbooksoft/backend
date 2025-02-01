package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type SignInCommand struct {
	Email    string
	Password string
}

type SignInHandlerResponse struct {
	AccessToken  string
	RefreshToken string
}

type SignInHandler struct {
	db *sql.DB
}

func NewSignInHandler(db *sql.DB) *SignInHandler {
	return &SignInHandler{db}
}

func (h *SignInHandler) Handle(ctx context.Context, cmd *SignInCommand) (*SignInHandlerResponse, error) {
	userEmail := cmd.Email
	password := cmd.Password

	var uPass string
	var uUserId string
	var uFirstName string
	var uLastName string
	err := h.db.QueryRow("SELECT id, first_name, last_name, password FROM users WHERE email = ? AND confirmed = 1", userEmail).Scan(&uUserId, &uFirstName, &uLastName, &uPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, gerr.NewBusinessLogicErr("userEmail or password invalid", gerr.AuthorizationUserNotFound)
		}

		log.Error().Err(err).Msgf("can not fetch user while auth %s", userEmail)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(uPass), []byte(password))

	if err != nil {
		return nil, gerr.NewBusinessLogicErr("userEmail or password invalid", gerr.AuthorizationUserNotFound)
	}

	at, rt, err := createTokens(h.db, uUserId, userEmail, uFirstName, uLastName)

	if err != nil {
		log.Error().Err(err).Msgf("can not create token when signing in %s", userEmail)
		return nil, fmt.Errorf("can not create token when signing in %w", err)
	}

	return &SignInHandlerResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}
