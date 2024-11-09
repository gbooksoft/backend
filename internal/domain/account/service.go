package account

import (
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/gerr"
)

type Service struct {
	db *sql.DB
}

func (s *Service) GetUserEmail(userId string) (string, error) {
	var uEmail sql.NullString
	err := s.db.QueryRow("SELECT email FROM users WHERE id = ?", userId).Scan(&uEmail)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", gerr.NewBusinessLogicErr("user not found", gerr.AccountUserNotFound, err)
		}
		return "", err
	}

	if !uEmail.Valid {
		return "", nil
	}

	return uEmail.String, nil
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}
