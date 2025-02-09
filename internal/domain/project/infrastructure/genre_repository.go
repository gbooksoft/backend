package infrastructure

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type MySqlGenreRepository struct {
	*repository.BaseMysqlRepository
}

func NewMySqlGenreRepository(db *sql.DB) *MySqlGenreRepository {
	return &MySqlGenreRepository{repository.NewBaseMysqlRepository(db)}
}

func (r *MySqlGenreRepository) Exists(tx *sql.Tx, genres []string) (bool, error) {
	if len(genres) == 0 {
		return false, nil
	}

	return true, nil
}
