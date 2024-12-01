package infrastructure

import (
	"database/sql"
	"encoding/json"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
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

	_, err = tx.Exec(query, character.MenuItemId.String(), character.Name, character.Description, string(bytes))

	if err != nil {
		return err
	}
	return nil
}
