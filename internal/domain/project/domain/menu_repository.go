package domain

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type MenuRepository interface {
	repository.BaseRepository
	AddMenu(tx *sql.Tx, menu *Menu) error
	FindByProjectId(tx *sql.Tx, userId UserId, projectId ProjectId) (*Menu, error)
	Save(tx *sql.Tx, userId UserId, projectId ProjectId, menu *Menu) error
}
