package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type CreateChapterCommand struct {
	Id          string
	ProjectId   string
	AuthorId    string
	Name        string
	Content     string
	ContentHtml string
}

type CreateChapterHandlerResponse struct {
	ProjectId string
}

type CreateChapterHandler struct {
	projectRepository domain.ProjectRepository
}

func NewCreateChapterHandler(projectRepository domain.ProjectRepository) *CreateChapterHandler {
	return &CreateChapterHandler{
		projectRepository: projectRepository,
	}
}

func (h *CreateChapterHandler) Handle(ctx context.Context, cmd *CreateChapterCommand) (*CreateChapterHandlerResponse, error) {
	ctx, err := h.projectRepository.BeginTransaction(context.TODO())
	if err != nil {
		return nil, err
	}
	defer h.projectRepository.RollbackTransaction(ctx)

	authorId := domain.AuthorId(uuid.MustParse(cmd.AuthorId))
	projectId := domain.ProjectId(uuid.MustParse(cmd.ProjectId))

	project, err := h.projectRepository.FindByAuthorIdAndProjectId(ctx, authorId, projectId)

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, gerr.NewBusinessLogicErr("project not found", gerr.NotFound)
	}

	newChapterId := domain.ChapterInternalId(uuid.MustParse(cmd.Id))

	exists, err := h.projectRepository.ChapterInternalIdExists(ctx, newChapterId)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, gerr.NewBusinessLogicErr("can not create chapter with this id", gerr.ProjectChapterAlreadyExists)
	}

	err = project.AddChapter(newChapterId, cmd.Name, cmd.Content, cmd.ContentHtml)

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
