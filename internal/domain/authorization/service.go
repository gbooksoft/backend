package authorization

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/internal/gerr"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"os"
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

func (s *Service) SignInUser(userEmail string, password string) (Token, error) {
	var uPass string
	var uUserId string
	var uFirstName string
	var uLastName string
	err := s.db.QueryRow("SELECT id, first_name, last_name, password FROM users WHERE email = ? AND confirmed = 1", userEmail).Scan(&uUserId, &uFirstName, &uLastName, &uPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Token{}, gerr.NewBusinessLogicErr("userEmail or password invalid", gerr.AuthorizationUserNotFound, err)
		}

		log.Error().Err(err).Msgf("can not fetch user while auth %s", userEmail)
		return Token{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(uPass), []byte(password))

	if err != nil {
		return Token{}, gerr.NewBusinessLogicErr("userEmail or password invalid", gerr.AuthorizationUserNotFound, err)
	}

	at, rt, err := createTokens(s.db, uUserId, uFirstName, uLastName)

	if err != nil {
		log.Error().Err(err).Msgf("can not create token when signing in %s", userEmail)
		return Token{}, fmt.Errorf("can not create token when signing in %w", err)
	}

	return Token{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (s *Service) Register(userEmail, password, firstName, lastName string) (string, error) {
	var exists sql.NullBool
	err := s.db.QueryRow("SELECT 1 FROM users WHERE email = ?", userEmail).Scan(&exists)

	if err == nil {
		return "", gerr.NewBusinessLogicErr("can not register user", gerr.AuthorizationUserAlreadyExists, err)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		log.Error().Err(err).Msgf("can not register user %s", userEmail)
		return "", err
	}

	// todo: check password length

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error().Err(err).Msgf("can hash password while register user %s", userEmail)
		return "", err
	}

	confirmationHash, err := s.generateConfirmationHash()

	if err != nil {
		log.Error().Err(err).Msgf("can not generate confirmation hash for user %s", userEmail)
		return "", err
	}

	tx, err := s.db.Begin()
	if err != nil {
		log.Error().Err(err).Msgf("can not start transaction")
		return "", err
	}
	defer tx.Rollback()

	userId := uuid.New().String()
	_, err = tx.Exec("INSERT INTO users (id, email, password, first_name, last_name, confirmation_hash ) VALUES (?, ?, ?, ?, ?, ?)", userId, userEmail, passwordHash, firstName, lastName, confirmationHash)

	if err != nil {
		log.Error().Err(err).Msgf("can not create user while register %s", userEmail)
		return "", err
	}

	link, err := s.generateConfirmationLink(userId, confirmationHash)

	if err != nil {
		log.Error().Err(err).Msgf("can not generate confirmation link for user %s", userEmail)
		return "", err
	}

	err = s.emailFacade.RegistrationConfirmationTx(tx, userEmail, email.RegistrationConfirmationData{
		FirstName:        firstName,
		LastName:         lastName,
		ConfirmationLink: link,
	})

	if err != nil {
		log.Error().Err(err).Msgf("can not add email to queue while register %s", userEmail)
		return "", err
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msgf("can not commit transaction while register %s", userEmail)
		return "", err
	}

	return userId, nil
}

func (s *Service) RefreshToken(rt string) (Token, error) {
	_, rtClaims, err := ParseToken(rt, true)

	if err != nil {
		log.Error().Err(err).Msgf("can not parse refresh token")
		return Token{}, err
	}

	if IsTokenExpired(rtClaims) {
		return Token{}, gerr.NewBusinessLogicErr("token expired", gerr.AuthorizationRefreshTokenExpired, err)
	}

	var at string
	err = s.db.QueryRow("SELECT access FROM auth_token WHERE refresh = ?", rt).Scan(&at)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Token{}, gerr.NewBusinessLogicErr("refresh not found", gerr.AuthorizationRefreshTokenNotFound, err)
		}

		log.Error().Err(err).Msgf("can not refresh token")
		return Token{}, err
	}

	_, atClaims, err := ParseToken(at, true)

	if err != nil {
		log.Error().Err(err).Msgf("can not parse refresh token")
		return Token{}, err
	}

	err = deleteByRefreshToken(s.db, rt)
	if err != nil {
		log.Error().Err(err).Msgf("can not delete refresh token")
		return Token{}, err
	}

	at, rt, err = createTokens(s.db, atClaims.Subject, atClaims.FirstName, atClaims.LastName)

	if err != nil {
		log.Error().Err(err).Msgf("can not create user when refreshing in %s", atClaims.Subject)
		return Token{}, fmt.Errorf("can not create user when refreshing in %w", err)
	}

	return Token{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (s *Service) Logout(rt string) error {
	_, rtClaims, err := ParseToken(rt, false)

	if err != nil {
		log.Error().Err(err).Msgf("can not parse refresh token")
		return err
	}

	if IsTokenExpired(rtClaims) {
		return gerr.NewBusinessLogicErr("token expired", gerr.AuthorizationRefreshTokenExpired, err)
	}

	err = deleteByRefreshToken(s.db, rt)
	if err != nil {
		log.Error().Err(err).Msgf("can not delete refresh token")
		return err
	}

	return nil
}

func (s *Service) generateConfirmationHash() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *Service) generateConfirmationLink(userId, hash string) (string, error) {
	u, err := url.Parse(os.Getenv("WEB_FRONT_CONFIRMATION_URL"))
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("user_id", userId)
	params.Add("token", hash)

	u.RawQuery = params.Encode()

	return u.String(), nil
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
			return gerr.NewBusinessLogicErr("confirmation not found", gerr.AuthorizationConfirmationNoFound, err)
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
