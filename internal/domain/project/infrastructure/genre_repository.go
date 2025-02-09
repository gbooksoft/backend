package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/pkg/repository"
	"strings"
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
	placeholders := make([]string, len(genres))
	values := make([]interface{}, len(genres))
	for i, genreID := range genres {
		placeholders[i] = "?"
		values[i] = genreID
	}

	query := fmt.Sprintf("SELECT count(id) as cnt FROM genres WHERE id IN (%s)", strings.Join(placeholders, ", "))

	var cnt int
	err := tx.QueryRow(query, values...).Scan(&cnt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return cnt == len(genres), nil
}
