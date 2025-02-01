package application

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/internal/domain/password"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UpdatePasswordCommand struct {
	UserId          string
	CurrentPassword string
	NewPassword     string
}

type UpdatePasswordHandlerResponse struct {
	ValidationError *password.PasswordValidationResult
}

type UpdatePasswordHandler struct {
	db          *sql.DB
	emailFacade *email.Facade
}

func NewUpdatePasswordHandler(db *sql.DB, emailFacade *email.Facade) *UpdatePasswordHandler {
	return &UpdatePasswordHandler{
		db:          db,
		emailFacade: emailFacade,
	}
}

func (h *UpdatePasswordHandler) Handle(ctx context.Context, cmd *UpdatePasswordCommand) (*UpdatePasswordHandlerResponse, error) {
	currentPassword := cmd.CurrentPassword
	newPassword := cmd.NewPassword
	userId := cmd.UserId

	validationResult := password.ValidatePassword(newPassword)
	if len(validationResult.Errors) > 0 {
		return &UpdatePasswordHandlerResponse{
			ValidationError: validationResult,
		}, nil
	}
	tx, err := h.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var uPassword string
	var uEmail string
	var uFirstName string
	var uLastName string
	err = tx.QueryRow("SELECT email, first_name, last_name, password FROM users WHERE id = ?", userId).Scan(&uEmail, &uFirstName, &uLastName, &uPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, gerr.NewBusinessLogicErr("user not found", gerr.AccountUserNotFound)
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(uPassword), []byte(currentPassword))

	if err != nil {
		return nil, gerr.NewBusinessLogicErr("current password is incorrect", gerr.AccountPasswordIsIncorrect)
	}

	if currentPassword == newPassword {
		err = h.emailFacade.PasswordChangedTx(tx, uEmail, email.PasswordChangedData{
			FirstName: uFirstName,
			LastName:  uLastName,
		})
		if err != nil {
			return nil, err
		}

		if err = tx.Commit(); err != nil {
			return nil, err
		}

		return nil, nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		log.Error().Err(err).Msgf("can hash password while changing it")
		return nil, err
	}

	_, err = tx.Exec("UPDATE users SET password = ? WHERE id = ?", passwordHash, userId)

	if err != nil {
		log.Error().Err(err).Msgf("can update password while sql update")
		return nil, err
	}

	err = h.emailFacade.PasswordChangedTx(tx, uEmail, email.PasswordChangedData{
		FirstName: uFirstName,
		LastName:  uLastName,
	})
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
