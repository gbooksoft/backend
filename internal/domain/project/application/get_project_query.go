package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type GetProjectQuery struct {
	AuthorId  string
	ProjectId string
}

type GetProjectHandlerResponse struct {
	Project readmodel.Project
}

type GetProjectHandler struct {
	projectReadModel readmodel.ProjectReadModel
}

func NewGetProjectHandler(projectReadModel readmodel.ProjectReadModel) *GetProjectHandler {
	return &GetProjectHandler{projectReadModel: projectReadModel}
}

func (h *GetProjectHandler) Handle(ctx context.Context, query *GetProjectQuery) (*GetProjectHandlerResponse, error) {
	project, err := h.projectReadModel.GetById(query.AuthorId, query.ProjectId)
	if err != nil {
		return nil, err
	}

	if project.Id == "" {
		return nil, nil
	}

	return &GetProjectHandlerResponse{Project: project}, nil
}
