package domain

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type GenreRepository interface {
	repository.BaseRepository
	Exists(tx *sql.Tx, genres []string) (bool, error)
}
