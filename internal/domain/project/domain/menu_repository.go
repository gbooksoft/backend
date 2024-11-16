package domain

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type MenuRepository interface {
	repository.BaseRepository
	AddMenu(tx *sql.Tx, menu []*MenuItem) error
}
