package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/google/uuid"
)

type RepublishProjectCommand struct {
	ProjectId string
	AuthorId  string
	Chapters  []string
}

type RepublishProjectHandlerResponse struct {
}

type RepublishProjectHandler struct {
	projectService domain.ProjectService
}

func NewRepublishProjectHandler(projectService domain.ProjectService) *RepublishProjectHandler {
	return &RepublishProjectHandler{
		projectService: projectService,
	}
}

func (h *RepublishProjectHandler) Handle(ctx context.Context, cmd *RepublishProjectCommand) (*RepublishProjectHandlerResponse, error) {
	authorId := domain.AuthorId(uuid.MustParse(cmd.AuthorId))
	projectId := domain.ProjectId(uuid.MustParse(cmd.ProjectId))

	var chaptersIds []domain.CharacterInternalId
	for _, ch := range cmd.Chapters {
		chaptersIds = append(chaptersIds, domain.CharacterInternalId(uuid.MustParse(ch)))
	}

	err := h.projectService.Republish(ctx, authorId, projectId, chaptersIds)

	if err != nil {
		return nil, err
	}

	return &RepublishProjectHandlerResponse{}, nil
}
