package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
	"github.com/google/uuid"
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
SELECT pmi.id, pmi.title, pmi.type, pmi.parent_id, pmi.is_selected, pmi.order
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
		var isSelected int8
		var menuItem readmodel.MenuItem
		if err := rows.Scan(&menuItem.Id, &menuItem.Title, &menuItem.ItemType, &menuItem.ParentId, &isSelected, &menuItem.Order); err != nil {
			return nil, err
		}
		if isSelected == 1 {
			menuItem.IsSelected = true
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

func (r *MySqlMenuReadModel) GetExpandedItems(userId string, projectId string) ([]string, error) {
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
			return []string{}, nil
		}

		log.Error().Err(err).Msg("Failed to get expanded items")
		return []string{}, err
	}

	var result []string
	err = json.Unmarshal([]byte(expandedItems), &result)

	if err != nil {
		log.Error().Err(err).Msg("Failed to get expanded items")
		return []string{}, err
	}

	return result, nil
}

func (r *MySqlMenuReadModel) ExpandItem(userId string, projectId string, menuItemId string) error {
	items, err := r.GetExpandedItems(userId, projectId)

	if err != nil {
		return err
	}

	shouldInsert := len(items) == 0

	items = append(items, menuItemId)

	bytes, err := json.Marshal(items)

	if err != nil {
		return err
	}

	var dbErr error
	if shouldInsert {
		_, dbErr = r.db.Exec("INSERT INTO project_menu_items_expanded (id, project_id, expanded_items) VALUES (?, ?, ?)", uuid.New().String(), projectId, string(bytes))
	} else {
		_, dbErr = r.db.Exec("UPDATE project_menu_items_expanded SET expanded_items = ? WHERE project_id = ?", string(bytes), projectId)
	}

	if dbErr != nil {
		return dbErr
	}

	return nil
}

func (r *MySqlMenuReadModel) CollapseItem(userId string, projectId string, menuItemId string) error {
	items, err := r.GetExpandedItems(userId, projectId)

	if err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	var newSlice []string
	for _, v := range items {
		if v != menuItemId {
			newSlice = append(newSlice, v)
		}
	}

	if len(newSlice) == 0 {
		_, err = r.db.Exec("DELETE FROM project_menu_items_expanded WHERE project_id = ?", projectId)

		if err != nil {
			return err
		}
	}

	bytes, err := json.Marshal(newSlice)

	if err != nil {
		return err
	}

	_, err = r.db.Exec("UPDATE project_menu_items_expanded SET expanded_items = ? WHERE project_id = ?", string(bytes), projectId)

	if err != nil {
		return err
	}

	return nil
}
