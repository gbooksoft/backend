package authorization

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/internal/gerr"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
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

	at, rt, err := createTokens(s.db, uUserId, uEmail)

	if err != nil {
		log.Error().Err(err).Msgf("can not create user when signing in %s", email)
		return Token{}, fmt.Errorf("can not create user when signing in %w", err)
	}

	return Token{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (s *Service) Register(email string, password string) (string, error) {
	var exists sql.NullBool
	err := s.db.QueryRow("SELECT 1 FROM users WHERE email = ?", email).Scan(&exists)

	if err == nil {
		return "", gerr.NewBusinessLogicErr("can not register user with email", err)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	// todo: check password length

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	userId := uuid.New().String()
	_, err = s.db.Exec("INSERT INTO users (id, email, password) VALUES (?, ?, ?)", userId, email, hash)

	if err != nil {
		log.Error().Err(err).Msgf("can not create user while register %s", email)
		return "", err
	}

	return userId, nil
}

func (s *Service) RefreshToken(rt string) (Token, error) {
	_, rtClaims, err := ParseToken(rt)

	if err != nil {
		log.Error().Err(err).Msgf("can not parse refresh token")
		return Token{}, gerr.NewBusinessLogicErr("can not parse refresh token", err)
	}

	if IsTokenExpired(rtClaims) {
		return Token{}, gerr.NewBusinessLogicErr("token expired", err)
	}

	var at string
	err = s.db.QueryRow("SELECT access FROM auth_token WHERE refresh = ?", rt).Scan(&at)

	if err != nil {
		log.Error().Err(err).Msgf("can not refresh token")
		return Token{}, gerr.NewBusinessLogicErr("can not refresh token", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return Token{}, gerr.NewBusinessLogicErr("refresh not found", err)
	}

	_, atClaims, err := ParseToken(at)

	if err != nil {
		return Token{}, gerr.NewBusinessLogicErr("can not refresh token, login again", err)
	}

	err = deleteByRefreshToken(s.db, rt)
	if err != nil {
		log.Error().Err(err).Msgf("can not delete refresh token")
		return Token{}, gerr.NewBusinessLogicErr("can not delete refresh token", err)
	}

	at, rt, err = createTokens(s.db, atClaims.Subject, atClaims.Email)

	if err != nil {
		log.Error().Err(err).Msgf("can not create user when refreshing in %s", atClaims.Email)
		return Token{}, fmt.Errorf("can not create user when refreshing in %w", err)
	}

	return Token{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (s *Service) Logout(rt string) error {
	_, rtClaims, err := ParseToken(rt)

	if err != nil {
		log.Error().Err(err).Msgf("can not parse refresh token")
		return gerr.NewBusinessLogicErr("can not parse refresh token", err)
	}

	if IsTokenExpired(rtClaims) {
		return gerr.NewBusinessLogicErr("token expired", err)
	}

	err = deleteByRefreshToken(s.db, rt)
	if err != nil {
		log.Error().Err(err).Msgf("can not delete refresh token")
		return gerr.NewBusinessLogicErr("can not delete refresh token", err)
	}

	return nil
}
