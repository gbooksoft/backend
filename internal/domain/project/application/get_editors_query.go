package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type GetEditorsQuery struct {
	UserId    string
	ProjectId string
}

type GetEditorsHandlerResponse struct {
	Editors []readmodel.Editor
}

type GetEditorsHandler struct {
	editorReadModel readmodel.EditorReadModel
}

func NewGetEditorsHandler(editorReadModel readmodel.EditorReadModel) *GetEditorsHandler {
	return &GetEditorsHandler{editorReadModel: editorReadModel}
}

func (h *GetEditorsHandler) Handle(ctx context.Context, query *GetEditorsQuery) (*GetEditorsHandlerResponse, error) {
	editors, err := h.editorReadModel.GetEditors(query.UserId, query.ProjectId)

	if err != nil {
		return nil, err
	}

	return &GetEditorsHandlerResponse{Editors: editors}, nil
}
