package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type DeleteChapterCommand struct {
	ProjectId         string
	AuthorId          string
	ChapterInternalId string
}

type DeleteChapterHandlerResponse struct {
}

type DeleteChapterHandler struct {
	projectRepository domain.ProjectRepository
}

func NewDeleteChapterHandler(projectRepository domain.ProjectRepository) *DeleteChapterHandler {
	return &DeleteChapterHandler{
		projectRepository: projectRepository,
	}
}

func (h *DeleteChapterHandler) Handle(ctx context.Context, cmd *DeleteChapterCommand) (*DeleteChapterHandlerResponse, error) {
	ctx, err := h.projectRepository.BeginTransaction(context.TODO())
	if err != nil {
		return nil, err
	}
	defer h.projectRepository.RollbackTransaction(ctx)

	authorId := domain.AuthorId(uuid.MustParse(cmd.AuthorId))
	projectId := domain.ProjectId(uuid.MustParse(cmd.ProjectId))
	chapterInternalId := domain.ChapterInternalId(uuid.MustParse(cmd.ChapterInternalId))

	project, err := h.projectRepository.FindByAuthorIdAndProjectId(ctx, authorId, projectId)

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, gerr.NewBusinessLogicErr("project not found", gerr.NotFound)
	}

	err = project.DeleteChapter(chapterInternalId)

	if err != nil {
		return nil, err
	}

	err = h.projectRepository.Save(ctx, project)

	if err != nil {
		return nil, err
	}

	if err := h.projectRepository.CommitTransaction(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
