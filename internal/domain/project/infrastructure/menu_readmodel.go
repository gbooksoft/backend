package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
	"github.com/rs/zerolog/log"
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

func (r *MySqlMenuReadModel) GetExpandedItems(userId string, projectId string) []string {
	query := `
SELECT pmi.expanded_items 
FROM project_menu_items_expanded pmi
JOIN projects p on p.id = pmi.project_id
WHERE pmi.project_id = ? AND p.user_id = ?
LIMIT 1
`
	var expandedItems string
	err := r.db.QueryRow(query, projectId, userId).Scan(&expandedItems)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}
		}

		log.Error().Err(err).Msg("Failed to get expanded items")
		return []string{}
	}

	var result []string
	err = json.Unmarshal([]byte(expandedItems), &result)

	if err != nil {
		log.Error().Err(err).Msg("Failed to get expanded items")
		return []string{}
	}

	return result
}
