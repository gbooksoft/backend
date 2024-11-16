package infrastructure

import (
	"database/sql"
	"errors"
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

func (r *MySqlMenuRepository) AddMenu(tx *sql.Tx, menu *domain.Menu, parentId *string) error {
	return r.addMenuItems(tx, menu.Items, parentId)
}

func (r *MySqlMenuRepository) addMenuItems(tx *sql.Tx, menuItems []*domain.MenuItem, parentId *string) error {
	for _, menuItem := range menuItems {
		itemId := menuItem.Id.String()
		_, err := tx.Exec(
			"INSERT INTO project_menu_items (id, project_id, title, type, parent_id) VALUES (?, ?, ?, ?, ?)",
			menuItem.Id.String(),
			menuItem.ProjectId.String(),
			menuItem.Title,
			string(menuItem.Type),
			parentId,
		)

		if len(menuItem.Children) > 0 {
			err = r.addMenuItems(tx, menuItem.Children, &itemId)

			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MySqlMenuRepository) FindByProjectId(tx *sql.Tx, projectId domain.ProjectId) (*domain.Menu, error) {
	query := `
SELECT id, project_id, title, type, parent_id FROM project_menu_items WHERE project_id = ?
`
	rows, err := tx.Query(query, projectId.String())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupedItems := make(map[string][]*domain.MenuItem)
	for rows.Next() {
		var item domain.MenuItem
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

		parentId := "ROOT"
		if itemParentId.Valid {
			parentId = itemParentId.String
		}
		groupedItems[parentId] = append(groupedItems[parentId], &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	for parentId, children := range groupedItems {
		for childI, _ := range children {
			curId := groupedItems[parentId][childI].Id

			if itsChildren, ok := groupedItems[curId.String()]; ok {
				groupedItems[parentId][childI].Children = itsChildren
			}
		}
	}

	if items, ok := groupedItems["ROOT"]; ok {
		return &domain.Menu{
			Items: items,
		}, nil
	}

	return nil, errors.New("errored structure of menu")
}
func (r *MySqlMenuRepository) Save(tx *sql.Tx, projectId domain.ProjectId, newMenu *domain.Menu) error {
	// now: 1,2,3,4,5,6
	// db: 4,5,6
	// diff now:db = 1,2,3
	// insert 1,2,3
	// diff db:now = null
	// do nothing
	// intersect db:now = 4,5,6
	// update, 4,5,6,

	// now:  3,4,5
	// db: 5,6,7
	// diff now:db = 3,4
	// insert 3,4
	// diff db:now = 6,7
	// delete 6,7
	// intersect db:now = 5
	// update 5
	// Шаг 1: Получаем текущие данные из базы
	_, err := tx.Exec("DELETE FROM project_menu_items WHERE project_id = ?", projectId.String())
	if err != nil {
		return err
	}

	err = r.AddMenu(tx, newMenu, nil)

	if err != nil {
		return err
	}

	return nil
}
