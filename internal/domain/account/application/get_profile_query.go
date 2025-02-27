package application

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/pkg/gerr"
)

type GetProfileQuery struct {
	UserId string
}

type GetProfileHandlerResponse struct {
	Email      string
	Username   string
	FirstName  string
	LastName   string
	AvatarPath string
}

type GetProfileHandler struct {
	db *sql.DB
}

func NewGetProfileHandler(db *sql.DB) *GetProfileHandler {
	return &GetProfileHandler{
		db: db,
	}
}

func (h *GetProfileHandler) Handle(ctx context.Context, query *GetProfileQuery) (*GetProfileHandlerResponse, error) {
	var uUsername sql.NullString
	var uEmail sql.NullString
	var uFirstName sql.NullString
	var uLastName sql.NullString
	var uAvatarPath sql.NullString
	err := h.db.QueryRow("SELECT email, username, first_name, last_name, avatar_path FROM users WHERE id = ?", query.UserId).Scan(&uEmail, &uUsername, &uFirstName, &uLastName, &uAvatarPath)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, gerr.NewBusinessLogicErr("user not found", gerr.AccountUserNotFound)
		}
		return nil, err
	}

	return &GetProfileHandlerResponse{
		Email:      uEmail.String,
		Username:   uUsername.String,
		FirstName:  uFirstName.String,
		LastName:   uLastName.String,
		AvatarPath: uAvatarPath.String,
	}, nil
}
