package infrastructure

import (
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type MySqlEditorReadModel struct {
	db *sql.DB
}

func NewMySqlEditorReadModel(db *sql.DB) *MySqlEditorReadModel {
	return &MySqlEditorReadModel{db: db}
}

func (r *MySqlEditorReadModel) GetEditors(userId, projectId string) ([]*readmodel.Editor, error) {
	query := `
SELECT pe.editor_state, pe.project_menu_id 
FROM project_editors pe
JOIN project_menu_items pmi ON pmi.id = pe.project_menu_id
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

	result := make([]*readmodel.Editor, 0)

	for rows.Next() {
		editor := &readmodel.Editor{}
		if err := rows.Scan(&editor.EditorState, &editor.MenuItemId); err != nil {
			return nil, err
		}

		result = append(result, editor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
