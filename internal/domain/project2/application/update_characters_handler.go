package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project2/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type UpdateCharacterCommand struct {
	ProjectId           string
	AuthorId            string
	CharacterInternalId string
	Name                string
	Description         *string
	Attributes          map[string]string
}

type UpdateCharacterHandlerResponse struct {
}

type UpdateCharacterHandler struct {
	projectRepository domain.ProjectRepository
}

func NewUpdateCharacterHandler(projectRepository domain.ProjectRepository) *UpdateCharacterHandler {
	return &UpdateCharacterHandler{
		projectRepository: projectRepository,
	}
}

func (h *UpdateCharacterHandler) Handle(ctx context.Context, cmd *UpdateCharacterCommand) (*UpdateCharacterHandlerResponse, error) {
	ctx, err := h.projectRepository.BeginTransaction(context.TODO())
	if err != nil {
		return nil, err
	}
	defer h.projectRepository.RollbackTransaction(ctx)

	authorId := domain.AuthorId(uuid.MustParse(cmd.AuthorId))
	projectId := domain.ProjectId(uuid.MustParse(cmd.ProjectId))
	characterInternalId := domain.ChapterInternalId(uuid.MustParse(cmd.CharacterInternalId))

	project, err := h.projectRepository.FindByAuthorIdAndProjectId(ctx, authorId, projectId)

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, gerr.NewBusinessLogicErr("project not found", gerr.NotFound)
	}

	err = project.UpdateCharacter(characterInternalId, cmd.Name, cmd.Description, cmd.Attributes)

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
