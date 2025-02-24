package infrastructure

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/internal/domain/project/infrastructure/adapter"
	"github.com/gbooksoft/backend/pkg/gerr"
)

type ProjectService struct {
	projectRepository domain.ProjectRepository
	bookAdapter       *adapter.BookAdapter
}

func NewProjectService(projectRepository domain.ProjectRepository, bookAdapter *adapter.BookAdapter) *ProjectService {
	return &ProjectService{
		projectRepository: projectRepository,
		bookAdapter:       bookAdapter,
	}
}

func (s *ProjectService) Publish(authorId domain.AuthorId, projectId domain.ProjectId, chaptersToPublish []domain.ChapterInternalId) error {
	ctx, err := s.projectRepository.BeginTransaction(context.TODO())
	if err != nil {
		return err
	}
	defer s.projectRepository.RollbackTransaction(ctx)

	project, err := s.projectRepository.FindByAuthorIdAndProjectId(ctx, authorId, projectId)

	if err != nil {
		return err
	}

	if project == nil {
		return gerr.NewNotFoundErr("project not found")
	}

	if project.BookId != nil {
		return gerr.NewBusinessLogicErr("project already published", gerr.ProjectAlreadyPublished)
	}

	bookId, err := s.bookAdapter.Publish(ctx, project, chaptersToPublish)

	if err != nil {
		return err
	}

	project.BookId = bookId

	err = s.projectRepository.Save(ctx, project)
	if err != nil {
		return err
	}

	if err := s.projectRepository.CommitTransaction(ctx); err != nil {
		return err
	}

	return nil
}
