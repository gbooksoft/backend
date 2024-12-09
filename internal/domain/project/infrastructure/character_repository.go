package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
	"github.com/google/uuid"
)

type MySqlCharacterRepository struct {
	*repository.BaseMysqlRepository
}

func NewMySqlCharacterRepository(db *sql.DB) *MySqlCharacterRepository {
	return &MySqlCharacterRepository{repository.NewBaseMysqlRepository(db)}
}

func (r *MySqlCharacterRepository) AddCharacter(tx *sql.Tx, character *domain.Character) error {
	query := `
INSERT INTO project_characters (project_menu_id, name, description, attributes) VALUES (?, ?, ?, ?)
`
	bytes, err := json.Marshal(character.Attributes)

	if err != nil {
		return err
	}

	desc := sql.NullString{
		String: "",
		Valid:  false,
	}
	if character.Description != nil {
		desc.Valid = true
		desc.String = *character.Description
	}

	_, err = tx.Exec(query, character.MenuItemId.String(), character.Name, desc, string(bytes))

	if err != nil {
		return err
	}
	return nil
}

func (r *MySqlCharacterRepository) FindCharacter(tx *sql.Tx, userId domain.UserId, projectId domain.ProjectId, menuItemId domain.MenuItemId) (*domain.Character, error) {
	query := `
SELECT pc.project_menu_id, pc.name, pc.description, pc.attributes
FROM project_characters pc 
JOIN project_menu_items pmi ON pc.project_menu_id  = pmi.id
JOIN projects p ON pmi.project_id = p.id
WHERE p.user_id = ? AND p.id = ? AND pc.project_menu_id = ?
`

	var char domain.Character
	var projectMenuItemId string
	var attrs sql.NullString
	var desc sql.NullString
	err := tx.QueryRow(query, userId.String(), projectId.String(), menuItemId.String()).Scan(&projectMenuItemId, &char.Name, &desc, &attrs)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	char.MenuItemId = domain.MenuItemId(uuid.MustParse(projectMenuItemId))

	if attrs.Valid {
		err = json.Unmarshal([]byte(attrs.String), &char.Attributes)

		if err != nil {
			return nil, err
		}
	}

	if desc.Valid {
		char.Description = &desc.String
	}

	return &char, nil
}
func (r *MySqlCharacterRepository) Save(tx *sql.Tx, character *domain.Character) error {
	query := `
UPDATE project_characters 
SET name = ?, description = ?, attributes = ? 
WHERE project_menu_id = ?
`
	attrs := sql.NullString{
		String: "",
		Valid:  false,
	}

	if len(character.Attributes) > 0 {
		data, err := json.Marshal(character.Attributes)
		if err != nil {
			return err
		}

		attrs.Valid = true
		attrs.String = string(data)
	}

	desc := sql.NullString{
		String: "",
		Valid:  false,
	}
	if character.Description != nil {
		desc.Valid = true
		desc.String = *character.Description
	}

	_, err := tx.Exec(query, character.Name, desc, attrs, character.MenuItemId.String())

	if err != nil {
		return err
	}

	return nil
}
