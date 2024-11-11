package account

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/internal/gerr"
	"github.com/gbooksoft/backend/pkg/token"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"os"
	"strings"
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
}

type UpdatedFields struct {
	Email     bool
	FirstName bool
	LastName  bool
}

func (s *Service) Update(userId string, data UpdateFields) (UpdatedFields, error) {
	var setClauses []string
	var args []interface{}
	argIndex := 1
	var result UpdatedFields

	if data.FirstName != nil {
		setClauses = append(setClauses, fmt.Sprintf("first_name = ?"))
		args = append(args, *data.FirstName)
		argIndex++
		result.FirstName = true
	}
	if data.LastName != nil {
		setClauses = append(setClauses, fmt.Sprintf("last_name = ?"))
		args = append(args, *data.LastName)
		argIndex++
		result.LastName = true
	}

	tx, err := s.db.Begin()
	if err != nil {
		return UpdatedFields{}, err
	}
	defer tx.Rollback()

	// обновляем обычные поля
	if len(setClauses) > 0 {
		query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(setClauses, ", "))
		args = append(args, userId)

		// Выполняем SQL запрос
		_, err = tx.Exec(query, args...)

		if err != nil {
			log.Fatal().Err(err).Msg("failed to update user")
			return UpdatedFields{}, err
		}
	}

	// составляем запрос на смену email
	if data.Email != nil {
		err = s.addChangeEmailRequest(tx, userId, *data.Email)

		if err != nil {
			log.Error().Err(err).Msg("failed to update users email")
			return UpdatedFields{}, err
		}

		result.Email = true
	}

	if err = tx.Commit(); err != nil {
		return UpdatedFields{}, err
	}

	return result, nil
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
			return gerr.NewBusinessLogicErr("token not found", gerr.AccountEmailChangeRequestNotFound, err)
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
		return gerr.NewBusinessLogicErr("token expired", gerr.AccountEmailChangeRequestExpired, nil)
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

func (s *Service) addChangeEmailRequest(tx *sql.Tx, userId, newEmail string) error {
	var exists sql.NullBool
	err := tx.QueryRow("SELECT 1 FROM users WHERE email = ?", newEmail).Scan(&exists)
	if err == nil {
		return gerr.NewBusinessLogicErr("user with this email already exists", gerr.AccountUserWithEmailAlreadyExists, nil)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	err = tx.QueryRow("SELECT 1 FROM email_change_requests WHERE new_email = ?", newEmail).Scan(&exists)

	if err == nil {
		return gerr.NewBusinessLogicErr("request on this email already exists", gerr.AccountEmailChangeRequestWithEmailAlreadyExists, nil)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	newToken, err := token.GenerateRandHash(64)

	if err != nil {
		return err
	}

	err = tx.QueryRow("SELECT 1 FROM email_change_requests WHERE user_id = ?", userId).Scan(&exists)

	// update and return
	if err == nil {
		_, err = tx.Exec("UPDATE email_change_requests SET new_email = ?, token = ?, created_at = NOW() WHERE user_id = ?", newEmail, newToken, userId)
		if err != nil {
			return err
		}
	}

	// insert and return
	if errors.Is(err, sql.ErrNoRows) {
		_, err = tx.Exec("INSERT INTO email_change_requests (user_id, new_email, token) VALUES (?, ?, ?)", userId, newEmail, newToken)
		if err != nil {
			return err
		}
	}

	var oldEmail sql.NullString
	var firstName sql.NullString
	var lastName sql.NullString
	err = tx.QueryRow("SELECT email, first_name, last_name FROM users WHERE id = ?", userId).Scan(&oldEmail, &firstName, &lastName)
	if err != nil {
		return err
	}

	confirmationLink, err := s.generateChangeEmailConfirmationLink(newToken)
	if err != nil {
		return err
	}
	err = s.emailFacade.EmailChangeRequestConfirmationTx(tx, newEmail, email.EmailChangeRequestConfirmationData{
		FirstName:        firstName.String,
		LastName:         lastName.String,
		ConfirmationLink: confirmationLink,
	})
	if err != nil {
		return err
	}

	err = s.emailFacade.EmailChangeRequestInformationTx(tx, oldEmail.String, email.EmailChangeRequestInformationData{
		FirstName: firstName.String,
		LastName:  lastName.String,
		NewEmail:  newEmail,
	})
	if err != nil {
		return err
	}

	return err

}

func (s *Service) generateChangeEmailConfirmationLink(hash string) (string, error) {
	u, err := url.Parse(os.Getenv("WEB_FRONT_EMAIL_CHANGE_REQUEST_CONFIRMATION_URL"))
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("token", hash)

	u.RawQuery = params.Encode()

	return u.String(), nil
}

func (s *Service) UpdatePassword(userId string, currentPassword string, newPassword string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var uPassword string
	var uEmail string
	var uFirstName string
	var uLastName string
	err = tx.QueryRow("SELECT email, first_name, last_name, password FROM users WHERE id = ?", userId).Scan(&uEmail, &uFirstName, &uLastName, &uPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gerr.NewBusinessLogicErr("user not found", gerr.AccountUserNotFound, nil)
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(uPassword), []byte(currentPassword))

	if err != nil {
		return gerr.NewBusinessLogicErr("current password is incorrect", gerr.AccountPasswordIsIncorrect, nil)
	}

	if currentPassword == newPassword {
		err = s.emailFacade.PasswordChangedTx(tx, uEmail, email.PasswordChangedData{
			FirstName: uFirstName,
			LastName:  uLastName,
		})
		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		return nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		log.Error().Err(err).Msgf("can hash password while changing it")
		return err
	}

	_, err = tx.Exec("UPDATE users SET password = ? WHERE id = ?", passwordHash, userId)

	if err != nil {
		log.Error().Err(err).Msgf("can update password while sql update")
		return err
	}

	err = s.emailFacade.PasswordChangedTx(tx, uEmail, email.PasswordChangedData{
		FirstName: uFirstName,
		LastName:  uLastName,
	})
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
