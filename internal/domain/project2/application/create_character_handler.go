package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project2/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type CreateCharacterCommand struct {
	Id        string
	ProjectId string
	AuthorId  string
	Name      string
}

type CreateCharacterHandlerResponse struct {
	ProjectId string
}

type CreateCharacterHandler struct {
	projectRepository domain.ProjectRepository
}

func NewCreateCharacterHandler(projectRepository domain.ProjectRepository) *CreateCharacterHandler {
	return &CreateCharacterHandler{
		projectRepository: projectRepository,
	}
}

func (h *CreateCharacterHandler) Handle(ctx context.Context, cmd *CreateCharacterCommand) (*CreateCharacterHandlerResponse, error) {
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

	err = project.AddCharacter(domain.CharacterInternalId(uuid.MustParse(cmd.Id)), cmd.Name, nil, nil)

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
