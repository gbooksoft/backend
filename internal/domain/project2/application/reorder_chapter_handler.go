package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project2/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type ReorderChapterCommand struct {
	ChapterInternalId string
	ProjectId         string
	AuthorId          string
	Order             int
}

type ReorderChapterHandlerResponse struct {
	ProjectId string
}

type ReorderChapterHandler struct {
	projectRepository domain.ProjectRepository
}

func NewReorderChapterHandler(projectRepository domain.ProjectRepository) *ReorderChapterHandler {
	return &ReorderChapterHandler{
		projectRepository: projectRepository,
	}
}

func (h *ReorderChapterHandler) Handle(ctx context.Context, cmd *ReorderChapterCommand) (*ReorderChapterHandlerResponse, error) {
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

	project.ReorderChapter(chapterInternalId, cmd.Order)

	err = h.projectRepository.Save(ctx, project)

	if err != nil {
		return nil, err
	}

	if err := h.projectRepository.CommitTransaction(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
