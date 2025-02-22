package domain

import (
	"context"
	"github.com/gbooksoft/backend/pkg/repository"
)

type ProjectRepository interface {
	repository.BaseRepository
	Add(ctx context.Context, project *Project) error
	Save(ctx context.Context, project *Project) error
	ExistsByUserIdAndTitle(ctx context.Context, authorId AuthorId, title string, ExceptProjectId *ProjectId) (bool, error)
	FindByAuthorIdAndProjectId(ctx context.Context, authorId AuthorId, projectId ProjectId) (*Project, error)
	ChapterInternalIdExists(ctx context.Context, chapterInternalId ChapterInternalId) (bool, error)
	CharacterInternalIdExists(ctx context.Context, characterInternalId CharacterInternalId) (bool, error)
	LocationInternalIdExists(ctx context.Context, locationInternalId LocationInternalId) (bool, error)
}

type GenreRepository interface {
	repository.BaseRepository
	Exists(ctx context.Context, genres []GenreId) (bool, error)
}
