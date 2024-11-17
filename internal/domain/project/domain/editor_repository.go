package domain

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type EditorRepository interface {
	repository.BaseRepository
	AddEditor(tx *sql.Tx, editor *Editor) error
	AddMany(tx *sql.Tx, editors []*Editor) error
}
