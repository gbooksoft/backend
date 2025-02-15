package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type CreateLocationCommand struct {
	Id        string
	ProjectId string
	AuthorId  string
	Name      string
}

type CreateLocationHandlerResponse struct {
	ProjectId string
}

type CreateLocationHandler struct {
	projectRepository domain.ProjectRepository
}

func NewCreateLocationHandler(projectRepository domain.ProjectRepository) *CreateLocationHandler {
	return &CreateLocationHandler{
		projectRepository: projectRepository,
	}
}

func (h *CreateLocationHandler) Handle(ctx context.Context, cmd *CreateLocationCommand) (*CreateLocationHandlerResponse, error) {
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

	newLocationId := domain.LocationInternalId(uuid.MustParse(cmd.Id))

	exists, err := h.projectRepository.LocationInternalIdExists(ctx, newLocationId)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, gerr.NewBusinessLogicErr("can not create location with this id", gerr.ProjectLocationAlreadyExists)
	}

	err = project.AddLocation(newLocationId, cmd.Name, nil)

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
