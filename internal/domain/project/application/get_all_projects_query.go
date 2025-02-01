package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type GetAllProjectsQuery struct {
	UserId string
}

type GetAllProjectsHandlerResponse struct {
	Projects []readmodel.Project
}

type GetAllProjectsHandler struct {
	projectReadModel readmodel.ProjectReadModel
}

func NewGetAllProjectsHandler(projectReadModel readmodel.ProjectReadModel) *GetAllProjectsHandler {
	return &GetAllProjectsHandler{projectReadModel: projectReadModel}
}

func (h *GetAllProjectsHandler) Handle(ctx context.Context, query *GetAllProjectsQuery) (*GetAllProjectsHandlerResponse, error) {
	projects, err := h.projectReadModel.GetAll(query.UserId)

	if err != nil {
		return nil, err
	}

	return &GetAllProjectsHandlerResponse{Projects: projects}, nil
}
