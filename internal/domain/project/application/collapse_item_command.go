package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type CollapseItemCommand struct {
	UserId     string
	ProjectId  string
	MenuItemId string
}

type CollapseItemHandlerResponse struct {
}

type CollapseItemHandler struct {
	menuReadModel readmodel.MenuReadModel
}

func NewCollapseItemHandler(menuReadModel readmodel.MenuReadModel) *CollapseItemHandler {
	return &CollapseItemHandler{menuReadModel: menuReadModel}
}
func (h *CollapseItemHandler) Handle(ctx context.Context, cmd *CollapseItemCommand) (*CollapseItemHandlerResponse, error) {
	err := h.menuReadModel.CollapseItem(cmd.UserId, cmd.ProjectId, cmd.MenuItemId)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
