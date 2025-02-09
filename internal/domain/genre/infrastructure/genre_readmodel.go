package infrastructure

import (
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/genre/readmodel"
)

type MySqlGenreReadModel struct {
	db *sql.DB
}

func NewMySqlGenreReadModel(db *sql.DB) *MySqlGenreReadModel {
	return &MySqlGenreReadModel{db: db}
}

func (r *MySqlGenreReadModel) GetAllByName(name string) ([]readmodel.Genre, error) {
	var err error
	var rows *sql.Rows

	if name == "" {
		query := `
SELECT id, name FROM genres ORDER BY name LIMIT 10
`
		rows, err = r.db.Query(query)
	} else {
		query := `
SELECT id, name FROM genres WHERE name LIKE ?
`
		rows, err = r.db.Query(query, "%"+name+"%")
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var genres []readmodel.Genre
	for rows.Next() {
		genre := readmodel.Genre{}

		err = rows.Scan(&genre.Id, &genre.Name)

		if err != nil {
			return nil, err
		}

		genres = append(genres, genre)
	}

	return genres, nil
}
