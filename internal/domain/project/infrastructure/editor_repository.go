package infrastructure

import (
	"database/sql"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
)

type MySqlEditorRepository struct {
	*repository.BaseMysqlRepository
}

func NewMySqlEditorRepository(db *sql.DB) *MySqlEditorRepository {
	return &MySqlEditorRepository{repository.NewBaseMysqlRepository(db)}
}

func (r *MySqlEditorRepository) AddEditor(tx *sql.Tx, editor *domain.Editor) error {
	query := `	
INSERT INTO project_editors (id, project_menu_id) VALUES (?, ?)
`
	_, err := tx.Exec(query, editor.Id.String(), editor.MenuItemId.String())

	if err != nil {
		return err
	}

	return nil
}

func (r *MySqlEditorRepository) AddMany(tx *sql.Tx, editors []*domain.Editor) error {
	for _, editor := range editors {
		if err := r.AddEditor(tx, editor); err != nil {
			return err
		}
	}
	return nil
}
