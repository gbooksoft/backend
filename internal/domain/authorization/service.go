package authorization

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/internal/gerr"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *sql.DB
}

type Token struct {
	AccessToken     string
	AccessTokenExp  int64
	RefreshToken    string
	RefreshTokenExp int64
}

func (s *Service) SignInUser(email string, password string) (Token, error) {
	var uPass string
	var uEmail string
	var uUserId string
	err := s.db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&uUserId, &uEmail, &uPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Token{}, gerr.NewBusinessLogicErr("email or password invalid", err)
		}

		log.Error().Err(err).Msgf("can not fetch user while auth %s", email)
		return Token{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(uPass), []byte(password))

	if err != nil {
		return Token{}, fmt.Errorf("email or password invalid %w", err)
	}

	at, atExp, rt, rtExp, err := CreateTokens(uUserId, uEmail)

	if err != nil {
		log.Error().Err(err).Msgf("can not create user when signing in %s", email)
		return Token{}, fmt.Errorf("can not create user when signing in %w", err)
	}

	return Token{
		AccessToken:     at,
		AccessTokenExp:  atExp,
		RefreshToken:    rt,
		RefreshTokenExp: rtExp,
	}, nil
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}
