package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/google/uuid"
)

type CanAccessCommand struct {
	UserId    string
	ProjectId string
}

type CanAccessHandlerResponse struct {
	CanAccess bool
}

type CanAccessHandler struct {
	projectRepository domain.ProjectRepository
}

func NewCanAccessHandler(projectRepository domain.ProjectRepository) *CanAccessHandler {
	return &CanAccessHandler{
		projectRepository: projectRepository,
	}
}

func (h *CanAccessHandler) Handle(ctx context.Context, cmd *CanAccessCommand) (*CanAccessHandlerResponse, error) {
	has, err := h.projectRepository.UserHasAccess(domain.UserId(uuid.MustParse(cmd.UserId)), domain.ProjectId(uuid.MustParse(cmd.ProjectId)))

	if err != nil {
		return nil, err
	}

	return &CanAccessHandlerResponse{CanAccess: has}, nil
}
