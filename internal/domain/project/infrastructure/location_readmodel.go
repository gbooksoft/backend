package infrastructure

import (
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type MySqlLocationReadModel struct {
	db *sql.DB
}

func NewMySqlLocationReadModel(db *sql.DB) *MySqlLocationReadModel {
	return &MySqlLocationReadModel{db: db}
}

func (r *MySqlLocationReadModel) GetLocations(userId, projectId string) ([]readmodel.Location, error) {
	query := `
SELECT pl.project_menu_id, pl.name, pl.description
FROM project_locations pl
JOIN project_menu_items pmi ON pmi.id = pl.project_menu_id
JOIN projects p ON pmi.project_id = p.id
WHERE p.user_id = ? AND p.id = ?
`
	rows, err := r.db.Query(query, userId, projectId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}
	defer rows.Close()

	result := make([]readmodel.Location, 0)

	for rows.Next() {
		var desc sql.NullString
		location := readmodel.Location{}
		if err := rows.Scan(&location.MenuItemId, &location.Name, &desc); err != nil {
			return nil, err
		}

		if desc.Valid {
			location.Description = desc.String
		}

		result = append(result, location)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
