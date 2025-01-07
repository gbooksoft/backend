package domain

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type ProjectRepository interface {
	repository.BaseRepository
	ExistsByUserIdAndTitle(tx *sql.Tx, userId UserId, title string, ExceptProjectId *ProjectId) (bool, error)
	ExistsByUserIdAndProjectId(tx *sql.Tx, userId UserId, projectId ProjectId) (bool, error)
	Add(tx *sql.Tx, project *Project) error
	UserHasAccess(userId UserId, projectId ProjectId) (bool, error)
	FindMetadataByUserIdAndProjectId(tx *sql.Tx, userId UserId, projectId ProjectId) (*Metadata, error)
	AddMetadata(tx *sql.Tx, metadata *Metadata) error
	SaveMetadata(tx *sql.Tx, metadata *Metadata) error
	FindByUserIdAndProjectId(tx *sql.Tx, userId UserId, projectId ProjectId) (*Project, error)
	Save(tx *sql.Tx, project *Project) error
}
