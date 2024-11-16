package domain

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type MenuRepository interface {
	repository.BaseRepository
	AddMenu(tx *sql.Tx, menu *Menu, parentId *string) error
	FindByProjectId(tx *sql.Tx, projectId ProjectId) (*Menu, error)
	Save(tx *sql.Tx, projectId ProjectId, menu *Menu) error
}
