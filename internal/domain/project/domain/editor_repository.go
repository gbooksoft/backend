package domain

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type EditorRepository interface {
	repository.BaseRepository
	AddEditor(tx *sql.Tx, editor *Editor) error
	AddMany(tx *sql.Tx, editors []*Editor) error
	FindEditor(tx *sql.Tx, User UserId, ProjectId ProjectId, id MenuItemId) (*Editor, error)
	Save(tx *sql.Tx, editor *Editor) error
}
