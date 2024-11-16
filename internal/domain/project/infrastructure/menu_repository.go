package infrastructure

import (
	"database/sql"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
)

type MySqlMenuRepository struct {
	*repository.BaseMysqlRepository
}

func NewMySqlMenuRepository(db *sql.DB) *MySqlMenuRepository {
	return &MySqlMenuRepository{repository.NewBaseMysqlRepository(db)}
}

func (r *MySqlMenuRepository) AddMenu(tx *sql.Tx, menu []*domain.MenuItem) error {
	for _, menuItem := range menu {
		_, err := tx.Exec(
			"INSERT INTO project_menu_items (id, project_id, title, type) VALUES (?, ?, ?, ?)",
			menuItem.Id.String(),
			menuItem.ProjectId.String(),
			menuItem.Title,
			string(menuItem.Type),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
