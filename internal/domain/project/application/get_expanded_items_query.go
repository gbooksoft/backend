package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type GetExpandedItemsQuery struct {
	UserId    string
	ProjectId string
}

type GetExpandedItemsHandlerResponse struct {
	ItemsIds []string
}
type GetExpandedItemsHandler struct {
	menuReadModel readmodel.MenuReadModel
}

func NewGetExpandedItemsHandler(menuReadModel readmodel.MenuReadModel) *GetExpandedItemsHandler {
	return &GetExpandedItemsHandler{menuReadModel: menuReadModel}
}

func (h *GetExpandedItemsHandler) Handle(ctx context.Context, query *GetExpandedItemsQuery) (*GetExpandedItemsHandlerResponse, error) {
	items, err := h.menuReadModel.GetExpandedItems(query.UserId, query.ProjectId)

	if err != nil {
		return nil, err
	}

	return &GetExpandedItemsHandlerResponse{ItemsIds: items}, nil
}
