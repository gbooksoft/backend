package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gbooksoft/backend/pkg/token"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"strings"
)

type UpdateFields struct {
	Email      *string
	FirstName  *string
	LastName   *string
	AvatarPath *string
}

type UpdateSettingsCommand struct {
	UserId       string
	UpdateFields UpdateFields
}

type updatedFields struct {
	Email      bool
	FirstName  bool
	LastName   bool
	AvatarPath bool
}

type UpdateSettingsHandlerResponse struct {
	UpdatedFields updatedFields
}

type UpdateSettingsHandler struct {
	db          *sql.DB
	emailFacade *email.Facade
}

func NewUpdateSettingsHandler(db *sql.DB, emailFacade *email.Facade) *UpdateSettingsHandler {
	return &UpdateSettingsHandler{
		db:          db,
		emailFacade: emailFacade,
	}
}

func (h *UpdateSettingsHandler) Handle(ctx context.Context, cmd *UpdateSettingsCommand) (*UpdateSettingsHandlerResponse, error) {
	userId := cmd.UserId
	var setClauses []string
	var args []interface{}
	argIndex := 1
	data := cmd.UpdateFields
	var result updatedFields

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
	if data.AvatarPath != nil {
		setClauses = append(setClauses, fmt.Sprintf("avatar_path = ?"))
		args = append(args, *data.AvatarPath)
		argIndex++
		result.AvatarPath = true
	}

	tx, err := h.db.Begin()
	if err != nil {
		return nil, err
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
			return nil, err
		}
	}

	// составляем запрос на смену email
	if data.Email != nil {
		err = h.addChangeEmailRequest(tx, userId, *data.Email)

		if err != nil {
			log.Error().Err(err).Msg("failed to update users email")
			return nil, err
		}

		result.Email = true
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &UpdateSettingsHandlerResponse{UpdatedFields: result}, nil
}

func (h *UpdateSettingsHandler) addChangeEmailRequest(tx *sql.Tx, userId, newEmail string) error {
	var exists sql.NullBool
	err := tx.QueryRow("SELECT 1 FROM users WHERE email = ?", newEmail).Scan(&exists)
	if err == nil {
		return gerr.NewBusinessLogicErr("user with this email already exists", gerr.AccountUserWithEmailAlreadyExists)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	err = tx.QueryRow("SELECT 1 FROM email_change_requests WHERE new_email = ?", newEmail).Scan(&exists)

	if err == nil {
		return gerr.NewBusinessLogicErr("request on this email already exists", gerr.AccountEmailChangeRequestWithEmailAlreadyExists)
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

	confirmationLink, err := h.generateChangeEmailConfirmationLink(newToken)
	if err != nil {
		return err
	}
	err = h.emailFacade.EmailChangeRequestConfirmationTx(tx, newEmail, email.EmailChangeRequestConfirmationData{
		FirstName:        firstName.String,
		LastName:         lastName.String,
		ConfirmationLink: confirmationLink,
	})
	if err != nil {
		return err
	}

	err = h.emailFacade.EmailChangeRequestInformationTx(tx, oldEmail.String, email.EmailChangeRequestInformationData{
		FirstName: firstName.String,
		LastName:  lastName.String,
		NewEmail:  newEmail,
	})
	if err != nil {
		return err
	}

	return err

}

func (h *UpdateSettingsHandler) generateChangeEmailConfirmationLink(hash string) (string, error) {
	u, err := url.Parse(os.Getenv("WEB_FRONT_EMAIL_CHANGE_REQUEST_CONFIRMATION_URL"))
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("token", hash)

	u.RawQuery = params.Encode()

	return u.String(), nil
}
