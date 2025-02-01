package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
	"github.com/gbooksoft/backend/pkg/gerr"
)

type GetMetadataQuery struct {
	UserId    string
	ProjectId string
}

type GetMetadataHandlerResponse struct {
	Metadata *readmodel.Metadata
}

type GetMetadataHandler struct {
	projectReadModel readmodel.ProjectReadModel
}

func NewGetMetadataHandler(projectReadModel readmodel.ProjectReadModel) *GetMetadataHandler {
	return &GetMetadataHandler{projectReadModel: projectReadModel}
}

func (h *GetMetadataHandler) Handle(ctx context.Context, query *GetMetadataQuery) (*GetMetadataHandlerResponse, error) {
	metadata, err := h.projectReadModel.GetMetadata(query.UserId, query.ProjectId)

	if err != nil {
		return nil, err
	}

	if metadata == nil {
		return nil, gerr.NewBusinessLogicErr("metadata not found", gerr.NotFound)
	}

	return &GetMetadataHandlerResponse{Metadata: metadata}, nil
}
