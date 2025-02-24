package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/google/uuid"
)

type PublishProjectCommand struct {
	ProjectId string
	AuthorId  string
	Chapters  []string
}

type PublishProjectHandlerResponse struct {
}

type PublishProjectHandler struct {
	projectService domain.ProjectService
}

func NewPublishProjectHandler(projectService domain.ProjectService) *PublishProjectHandler {
	return &PublishProjectHandler{
		projectService: projectService,
	}
}

func (h *PublishProjectHandler) Handle(ctx context.Context, cmd *PublishProjectCommand) (*PublishProjectHandlerResponse, error) {
	authorId := domain.AuthorId(uuid.MustParse(cmd.AuthorId))
	projectId := domain.ProjectId(uuid.MustParse(cmd.ProjectId))

	var chaptersIds []domain.CharacterInternalId
	for _, ch := range cmd.Chapters {
		chaptersIds = append(chaptersIds, domain.CharacterInternalId(uuid.MustParse(ch)))
	}

	err := h.projectService.Publish(authorId, projectId, chaptersIds)

	if err != nil {
		return nil, err
	}

	return &PublishProjectHandlerResponse{}, nil
}
