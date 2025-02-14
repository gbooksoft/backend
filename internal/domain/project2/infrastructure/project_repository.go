package infrastructure

import (
	"context"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project2/domain"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/gbooksoft/backend/pkg/repository"
	"gorm.io/gorm"
)

type MySqlProjectRepository struct {
	*repository.BaseGormRepository
}

func NewMySqlProjectRepository() *MySqlProjectRepository {
	return &MySqlProjectRepository{
		repository.NewBaseGormRepository(database.GetGormDB()),
	}
}

func (r *MySqlProjectRepository) ExistsByUserIdAndTitle(ctx context.Context, authorId domain.AuthorId, title string, ExceptProjectId *domain.ProjectId) (bool, error) {
	tx, ok := repository.FetchGormTransactionFromContext(ctx)

	db := r.Db
	if ok {
		db = tx
	}

	project := &domain.Project{}

	if ExceptProjectId != nil {
		if err := db.Table("projects2").First(project, "title = ? AND author_id = ? AND id != ?", title, authorId.String(), ExceptProjectId.String()).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
	} else {
		if err := db.Table("projects2").First(project, "title = ? AND author_id = ?", title, authorId.String()).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
	}

	if project.Id.String() == "" {
		return false, nil
	}

	return true, nil
}

func (r *MySqlProjectRepository) Add(ctx context.Context, project *domain.Project) error {
	tx, ok := repository.FetchGormTransactionFromContext(ctx)

	db := r.Db
	if ok {
		db = tx
	}

	result := db.Table("projects2").Create(project)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func (r *MySqlProjectRepository) Save(ctx context.Context, project *domain.Project) error {
	tx, ok := repository.FetchGormTransactionFromContext(ctx)

	db := r.Db
	if ok {
		db = tx
	}

	result := db.Save(project)

	if err := result.Error; err != nil {
		return err
	}

	if err := db.Model(project).Association("Genres").Replace(project.Genres); err != nil {
		return err
	}

	return nil
}

func (r *MySqlProjectRepository) FindByAuthorIdAndProjectId(ctx context.Context, authorId domain.AuthorId, projectId domain.ProjectId) (*domain.Project, error) {
	tx, ok := repository.FetchGormTransactionFromContext(ctx)

	db := r.Db
	if ok {
		db = tx
	}

	project := &domain.Project{}

	if err := db.Table("projects2").
		Preload("Genres").
		Preload("Chapters").
		Preload("Locations").
		Preload("Characters").
		Preload("Metadata").
		First(project, "author_id = ? AND id = ?", authorId.String(), projectId.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return project, nil
}
