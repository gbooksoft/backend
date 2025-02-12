package infrastructure

import (
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
	"github.com/google/uuid"
)

type MySqlLocationRepository struct {
	*repository.BaseMysqlRepository
}

func NewMySqlLocationRepository(db *sql.DB) *MySqlLocationRepository {
	return &MySqlLocationRepository{repository.NewBaseMysqlRepository(db)}
}

func (r *MySqlLocationRepository) AddLocation(tx *sql.Tx, location *domain.Location) error {
	query := `
INSERT INTO project_locations (id, project_menu_id, name, description) VALUES (?, ?, ?, ?)
`
	desc := sql.NullString{
		String: "",
		Valid:  false,
	}
	if location.Description != nil {
		desc.Valid = true
		desc.String = *location.Description
	}

	_, err := tx.Exec(query, r.GetNextId().String(), location.MenuItemId.String(), location.Name, desc)

	if err != nil {
		return err
	}
	return nil
}

func (r *MySqlLocationRepository) FindLocation(tx *sql.Tx, userId domain.UserId, projectId domain.ProjectId, menuItemId domain.MenuItemId) (*domain.Location, error) {
	query := `
SELECT pl.project_menu_id, pl.name, pl.description
FROM project_locations pl 
JOIN project_menu_items pmi ON pl.project_menu_id  = pmi.id
JOIN projects p ON pmi.project_id = p.id
WHERE p.user_id = ? AND p.id = ? AND pl.project_menu_id = ?
`

	var location domain.Location
	var projectMenuItemId string
	var desc sql.NullString
	err := tx.QueryRow(query, userId.String(), projectId.String(), menuItemId.String()).Scan(&projectMenuItemId, &location.Name, &desc)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	location.MenuItemId = domain.MenuItemId(uuid.MustParse(projectMenuItemId))

	if desc.Valid {
		location.Description = &desc.String
	}

	return &location, nil
}
func (r *MySqlLocationRepository) Save(tx *sql.Tx, location *domain.Location) error {
	query := `
UPDATE project_locations
SET name = ?, description = ? 
WHERE project_menu_id = ?
`

	desc := sql.NullString{
		String: "",
		Valid:  false,
	}
	if location.Description != nil {
		desc.Valid = true
		desc.String = *location.Description
	}

	_, err := tx.Exec(query, location.Name, desc, location.MenuItemId.String())

	if err != nil {
		return err
	}

	return nil
}
