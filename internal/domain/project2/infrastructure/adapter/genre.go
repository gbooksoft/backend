package adapter

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/internal/domain/project2/domain"
	"github.com/gbooksoft/backend/pkg/repository"
	"strings"
)

type GenreAdapter struct {
	*repository.BaseMysqlRepository
}

func NewGenreAdapter(db *sql.DB) *GenreAdapter {
	return &GenreAdapter{repository.NewBaseMysqlRepository(db)}
}

func (g *GenreAdapter) Exists(tx *sql.Tx, genres []domain.GenreId) (bool, error) {
	if len(genres) == 0 {
		return false, nil
	}
	placeholders := make([]string, len(genres))
	values := make([]interface{}, len(genres))
	for i, genreId := range genres {
		placeholders[i] = "?"
		values[i] = genreId.String()
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
