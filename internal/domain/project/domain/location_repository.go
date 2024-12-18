package domain

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type LocationRepository interface {
	repository.BaseRepository
	AddLocation(tx *sql.Tx, character *Location) error
	FindLocation(tx *sql.Tx, userId UserId, projectId ProjectId, menuItemId MenuItemId) (*Location, error)
	Save(tx *sql.Tx, character *Location) error
}
