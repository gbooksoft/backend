package infrastructure

import (
	"database/sql"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
	"github.com/google/uuid"
)

type MySqlEditorRepository struct {
	*repository.BaseMysqlRepository
}

func NewMySqlEditorRepository(db *sql.DB) *MySqlEditorRepository {
	return &MySqlEditorRepository{repository.NewBaseMysqlRepository(db)}
}

func (r *MySqlEditorRepository) AddEditor(tx *sql.Tx, editor *domain.Editor) error {
	query := `	
INSERT INTO project_editors (project_menu_id) VALUES (?)
`
	_, err := tx.Exec(query, editor.MenuItemId.String())

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

func (r *MySqlEditorRepository) FindEditor(tx *sql.Tx, userId domain.UserId, projectId domain.ProjectId, menuItemId domain.MenuItemId) (*domain.Editor, error) {
	query := `
SELECT pe.project_menu_id, pe.editor_state 
FROM project_editors pe
JOIN project_menu_items pmi ON pmi.id = pe.project_menu_id
JOIN projects p ON p.id = pmi.project_id
WHERE pe.project_menu_id = ? AND p.user_id = ? AND pmi.project_id = ?
`
	var projectMenuItemId string
	var editorState sql.NullString
	err := tx.QueryRow(query, menuItemId.String(), userId.String(), projectId.String()).Scan(&projectMenuItemId, &editorState)

	if err != nil {
		return nil, err
	}

	editor := &domain.Editor{}
	editor.MenuItemId = domain.MenuItemId(uuid.MustParse(projectMenuItemId))

	if editorState.Valid {
		stV := editorState.String
		editor.EditorState = &stV
	}

	return editor, nil
}

func (r *MySqlEditorRepository) Save(tx *sql.Tx, editor *domain.Editor) error {
	query := `
UPDATE project_editors SET editor_state = ?
WHERE project_menu_id = ?
`
	_, err := tx.Exec(query, editor.EditorState, editor.MenuItemId.String())

	if err != nil {
		return err
	}

	return nil
}
