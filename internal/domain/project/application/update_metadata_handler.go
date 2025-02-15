package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type UpdateMetadataCommand struct {
	ProjectId string
	AuthorId  string
	Author    string
	Keywords  string
	Language  string
	Isbn      string
}

type UpdateMetadataHandlerResponse struct {
}

type UpdateMetadataHandler struct {
	projectRepository domain.ProjectRepository
}

func NewUpdateMetadataHandler(projectRepository domain.ProjectRepository) *UpdateMetadataHandler {
	return &UpdateMetadataHandler{
		projectRepository: projectRepository,
	}
}

func (h *UpdateMetadataHandler) Handle(ctx context.Context, cmd *UpdateMetadataCommand) (*UpdateMetadataHandlerResponse, error) {
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

	project.UpdateMetadata(
		cmd.Author,
		cmd.Keywords,
		cmd.Language,
		cmd.Isbn,
	)

	err = h.projectRepository.Save(ctx, project)

	if err != nil {
		return nil, err
	}

	if err := h.projectRepository.CommitTransaction(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
