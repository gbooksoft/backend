package infrastructure

import (
	"database/sql"
	"fmt"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
	"github.com/google/uuid"
)

type MySqlMenuRepository struct {
	*repository.BaseMysqlRepository
}

func NewMySqlMenuRepository(db *sql.DB) *MySqlMenuRepository {
	return &MySqlMenuRepository{repository.NewBaseMysqlRepository(db)}
}

func (r *MySqlMenuRepository) AddMenu(tx *sql.Tx, menu *domain.Menu) error {
	for _, menuItem := range menu.Items {
		err := r.addMenuItem(tx, menuItem)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MySqlMenuRepository) FindByProjectId(tx *sql.Tx, userId domain.UserId, projectId domain.ProjectId) (*domain.Menu, error) {
	query := `
SELECT pmi.id, pmi.project_id, pmi.title, pmi.type, pmi.parent_id
FROM project_menu_items pmi 
JOIN projects p ON pmi.project_id = p.id
WHERE pmi.project_id = ? AND p.user_id = ?
`
	rows, err := tx.Query(query, projectId.String(), userId.String())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itemsResult []*domain.MenuItem
	for rows.Next() {
		item := &domain.MenuItem{}
		var itemId string
		var itemProjectId string
		var itemTitle string
		var itemType string
		var itemParentId sql.NullString

		err := rows.Scan(
			&itemId,
			&itemProjectId,
			&itemTitle,
			&itemType,
			&itemParentId,
		)

		if err != nil {
			return nil, err
		}

		item.Id = domain.MenuItemId(uuid.MustParse(itemId))
		item.ProjectId = domain.ProjectId(uuid.MustParse(itemProjectId))
		item.Title = itemTitle
		item.Type = domain.MenuItemType(itemType)

		if itemParentId.Valid {
			pId := domain.MenuItemId(uuid.MustParse(itemParentId.String))
			item.ParentId = &pId
		} else {
			item.ParentId = nil
		}

		itemsResult = append(itemsResult, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &domain.Menu{
		Items: itemsResult,
	}, nil
}

func (r *MySqlMenuRepository) Save(tx *sql.Tx, userId domain.UserId, projectId domain.ProjectId, newMenu *domain.Menu) error {
	currDbMenu, err := r.FindByProjectId(tx, userId, projectId)

	if err != nil {
		return err
	}

	existingItems := make(map[string]*domain.MenuItem)
	for _, item := range currDbMenu.Items {
		curr := item
		existingItems[item.Id.String()] = curr
	}

	newItemsMap := make(map[string]*domain.MenuItem)
	for _, newItem := range newMenu.Items {
		newItemsMap[newItem.Id.String()] = newItem
	}

	for existingId := range existingItems {
		if _, found := newItemsMap[existingId]; !found {
			// Элемент есть в базе, но его нет в новых данных — удалить
			_, err = tx.Exec("DELETE FROM project_menu_items WHERE id = ?", existingId)
			if err != nil {
				return fmt.Errorf("error deleting item with id %s: %w", existingId, err)
			}
		}
	}

	for _, newItem := range newMenu.Items {
		if existingItem, found := existingItems[newItem.Id.String()]; found {
			if existingItem.Title != newItem.Title {
				_, err := tx.Exec(
					"UPDATE project_menu_items SET title = ? WHERE id = ?",
					newItem.Title, newItem.Id.String(),
				)
				if err != nil {
					return fmt.Errorf("error updating item with id %s: %w", newItem.Id, err)
				}
			}
		} else {
			err = r.addMenuItem(tx, newItem)

			if err != nil {
				return fmt.Errorf("error inserting item with id %s: %w", newItem.Id, err)
			}
		}
	}

	return nil
}

func (r *MySqlMenuRepository) addMenuItem(tx *sql.Tx, menuItem *domain.MenuItem) error {
	var parentId sql.NullString
	if menuItem.ParentId != nil {
		parentId = sql.NullString{
			String: menuItem.ParentId.String(),
			Valid:  true,
		}
	} else {
		parentId = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	_, err := tx.Exec(
		"INSERT INTO project_menu_items (id, project_id, title, type, parent_id) VALUES (?, ?, ?, ?, ?)",
		menuItem.Id.String(),
		menuItem.ProjectId.String(),
		menuItem.Title,
		string(menuItem.Type),
		parentId,
	)

	if err != nil {
		return err
	}

	return nil
}
