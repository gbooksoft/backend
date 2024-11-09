package account

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/internal/gerr"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetUser(userId string) (string, string, string, error) {
	var uEmail sql.NullString
	var uFirstName sql.NullString
	var uLastName sql.NullString
	err := s.db.QueryRow("SELECT email, first_name, last_name FROM users WHERE id = ?", userId).Scan(&uEmail, &uFirstName, &uLastName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", "", gerr.NewBusinessLogicErr("user not found", gerr.AccountUserNotFound, err)
		}
		return "", "", "", err
	}

	return uEmail.String, uFirstName.String, uLastName.String, nil
}

type UpdateFields struct {
	Email     *string
	FirstName *string
	LastName  *string
	Password  *string
}

func (s *Service) Update(userId string, data UpdateFields) error {
	var setClauses []string
	var args []interface{}
	argIndex := 1

	if data.Email != nil {
		if s.existsUserWithEmail(*data.Email) {
			return gerr.NewBusinessLogicErr("can not set this email", gerr.AccountUserAlreadyExists, nil)
		}
		setClauses = append(setClauses, fmt.Sprintf("email = ?"))
		args = append(args, *data.Email)
		argIndex++
	}

	if data.FirstName != nil {
		setClauses = append(setClauses, fmt.Sprintf("first_name = ?"))
		args = append(args, *data.FirstName)
		argIndex++
	}
	if data.LastName != nil {
		setClauses = append(setClauses, fmt.Sprintf("last_name = ?"))
		args = append(args, *data.LastName)
		argIndex++
	}
	if data.Password != nil {
		// В реальном приложении пароль нужно шифровать перед сохранением
		setClauses = append(setClauses, fmt.Sprintf("password = ?"))
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)

		if err != nil {
			return err
		}

		args = append(args, passwordHash)
		argIndex++
	}

	if len(setClauses) == 0 {
		return gerr.NewBusinessLogicErr("nothing to update", gerr.AccountNothingToUpdate, nil)
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	args = append(args, userId)

	// Выполняем SQL запрос
	_, err := s.db.Exec(query, args...)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to update user")
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
