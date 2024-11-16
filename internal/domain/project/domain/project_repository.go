package domain

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type ProjectRepository interface {
	repository.BaseRepository
	ExistsByUserIdAndTitle(tx *sql.Tx, userId UserId, title string) (bool, error)
	Add(tx *sql.Tx, project *Project) error
}
