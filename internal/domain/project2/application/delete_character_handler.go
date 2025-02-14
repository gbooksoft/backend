package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project2/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type DeleteCharacterCommand struct {
	ProjectId           string
	AuthorId            string
	CharacterInternalId string
}

type DeleteCharacterHandlerResponse struct {
}

type DeleteCharacterHandler struct {
	projectRepository domain.ProjectRepository
}

func NewDeleteCharacterHandler(projectRepository domain.ProjectRepository) *DeleteCharacterHandler {
	return &DeleteCharacterHandler{
		projectRepository: projectRepository,
	}
}

func (h *DeleteCharacterHandler) Handle(ctx context.Context, cmd *DeleteCharacterCommand) (*DeleteCharacterHandlerResponse, error) {
	ctx, err := h.projectRepository.BeginTransaction(context.TODO())
	if err != nil {
		return nil, err
	}
	defer h.projectRepository.RollbackTransaction(ctx)

	authorId := domain.AuthorId(uuid.MustParse(cmd.AuthorId))
	projectId := domain.ProjectId(uuid.MustParse(cmd.ProjectId))
	characterInternalId := domain.CharacterInternalId(uuid.MustParse(cmd.CharacterInternalId))

	project, err := h.projectRepository.FindByAuthorIdAndProjectId(ctx, authorId, projectId)

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, gerr.NewBusinessLogicErr("project not found", gerr.NotFound)
	}

	err = project.DeleteCharacter(characterInternalId)

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
