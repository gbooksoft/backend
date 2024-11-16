package infrastructure

import (
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type MySqlMenuReadModel struct {
	db *sql.DB
}

func NewMySqlMenuReadModel(db *sql.DB) *MySqlMenuReadModel {
	return &MySqlMenuReadModel{db: db}
}

func (r *MySqlMenuReadModel) GetByUserIdAndProjectId(userId, projectId string) (*readmodel.Menu, error) {
	query := `
SELECT pmi.id, pmi.title, pmi.type, pmi.parent_id
FROM project_menu_items pmi
JOIN projects p on p.id = pmi.project_id
WHERE p.user_id = ? AND pmi.project_id = ? 
ORDER BY FIELD(pmi.type, 'chapter', 'character', 'location'), pmi.created_at 
`
	rows, err := r.db.Query(query, userId, projectId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}
	defer rows.Close()

	var menuItems []readmodel.MenuItem

	for rows.Next() {
		var menuItem readmodel.MenuItem
		if err := rows.Scan(&menuItem.Id, &menuItem.Title, &menuItem.ItemType, &menuItem.ParentId); err != nil {
			return nil, err
		}
		menuItems = append(menuItems, menuItem)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &readmodel.Menu{
		Items: menuItems,
	}, nil
}
