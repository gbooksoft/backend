package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type GetPublicationQuery struct {
	AuthorId  string
	ProjectId string
}

type GetPublicationHandlerResponse struct {
	Publication readmodel.Publication
}

type GetPublicationHandler struct {
	projectReadModel readmodel.ProjectReadModel
}

func NewGetPublicationHandler(projectReadModel readmodel.ProjectReadModel) *GetPublicationHandler {
	return &GetPublicationHandler{projectReadModel: projectReadModel}
}

func (h *GetPublicationHandler) Handle(ctx context.Context, query *GetPublicationQuery) (*GetPublicationHandlerResponse, error) {
	p, err := h.projectReadModel.GetPublication(query.AuthorId, query.ProjectId)
	if err != nil {
		return nil, err
	}

	if len(p.Chapters) == 0 {
		return nil, nil
	}

	return &GetPublicationHandlerResponse{
		Publication: p,
	}, nil
}
