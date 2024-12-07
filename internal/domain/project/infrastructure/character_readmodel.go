package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type MySqlCharacterReadModel struct {
	db *sql.DB
}

func NewMySqlCharacterReadModel(db *sql.DB) *MySqlCharacterReadModel {
	return &MySqlCharacterReadModel{db: db}
}

func (r *MySqlCharacterReadModel) GetCharacters(userId, projectId string) ([]readmodel.Character, error) {
	query := `
SELECT pc.project_menu_id, pc.name, pc.description, pc.attributes
FROM project_characters pc
JOIN project_menu_items pmi ON pmi.id = pc.project_menu_id
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

	result := make([]readmodel.Character, 0)

	for rows.Next() {
		var attrs sql.NullString
		character := readmodel.Character{}
		if err := rows.Scan(&character.MenuItemId, &character.Name, &character.Description, &attrs); err != nil {
			return nil, err
		}

		if attrs.Valid {
			err := json.Unmarshal([]byte(attrs.String), &character.Attributes)
			if err != nil {
				return nil, err
			}
		}

		result = append(result, character)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
