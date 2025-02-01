package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type GetRecentProjectsQuery struct {
	UserId string
	Limit  int
}

type GetRecentProjectsHandlerResponse struct {
	Projects []readmodel.Project
}

type GetRecentProjectsHandler struct {
	projectReadModel readmodel.ProjectReadModel
}

func NewGetRecentProjectsQueryHandler(projectReadModel readmodel.ProjectReadModel) *GetRecentProjectsHandler {
	return &GetRecentProjectsHandler{projectReadModel: projectReadModel}
}

func (h *GetRecentProjectsHandler) Handle(ctx context.Context, query *GetRecentProjectsQuery) (*GetRecentProjectsHandlerResponse, error) {
	projects, err := h.projectReadModel.GetRecent(query.UserId, query.Limit)

	if err != nil {
		return nil, err
	}

	return &GetRecentProjectsHandlerResponse{Projects: projects}, nil
}
