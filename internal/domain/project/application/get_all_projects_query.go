package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type GetAllProjectsQuery struct {
	UserId string
}

type GetAllProjectsQueryHandlerResponse struct {
	Projects []readmodel.Project
}

type GetAllProjectsQueryHandler struct {
	projectReadModel readmodel.ProjectReadModel
}

func NewGetAllProjectsQueryHandler(projectReadModel readmodel.ProjectReadModel) *GetAllProjectsQueryHandler {
	return &GetAllProjectsQueryHandler{projectReadModel: projectReadModel}
}

func (h *GetAllProjectsQueryHandler) Handle(ctx context.Context, query *GetAllProjectsQuery) (*GetAllProjectsQueryHandlerResponse, error) {
	projects, err := h.projectReadModel.GetAll(query.UserId)

	if err != nil {
		return nil, err
	}

	return &GetAllProjectsQueryHandlerResponse{Projects: projects}, nil
}
