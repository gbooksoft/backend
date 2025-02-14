package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project2/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type UpdateLocationCommand struct {
	ProjectId          string
	AuthorId           string
	LocationInternalId string
	Name               string
	Description        *string
}

type UpdateLocationHandlerResponse struct {
}

type UpdateLocationHandler struct {
	projectRepository domain.ProjectRepository
}

func NewUpdateLocationHandler(projectRepository domain.ProjectRepository) *UpdateLocationHandler {
	return &UpdateLocationHandler{
		projectRepository: projectRepository,
	}
}

func (h *UpdateLocationHandler) Handle(ctx context.Context, cmd *UpdateLocationCommand) (*UpdateLocationHandlerResponse, error) {
	ctx, err := h.projectRepository.BeginTransaction(context.TODO())
	if err != nil {
		return nil, err
	}
	defer h.projectRepository.RollbackTransaction(ctx)

	authorId := domain.AuthorId(uuid.MustParse(cmd.AuthorId))
	projectId := domain.ProjectId(uuid.MustParse(cmd.ProjectId))
	locationInternalId := domain.LocationInternalId(uuid.MustParse(cmd.LocationInternalId))

	project, err := h.projectRepository.FindByAuthorIdAndProjectId(ctx, authorId, projectId)

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, gerr.NewBusinessLogicErr("project not found", gerr.NotFound)
	}

	err = project.UpdateLocation(locationInternalId, cmd.Name, cmd.Description)

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
