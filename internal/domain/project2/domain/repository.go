package domain

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type ProjectRepository interface {
	repository.BaseRepository
	Add(tx *sql.Tx, project *Project) error
	Save(tx *sql.Tx, project *Project) error
	ExistsByUserIdAndTitle(tx *sql.Tx, authorId AuthorId, title string) (bool, error)
}

type GenreRepository interface {
	repository.BaseRepository
	Exists(tx *sql.Tx, genres []GenreId) (bool, error)
}
