package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type ExpandItemCommand struct {
	UserId     string
	ProjectId  string
	MenuItemId string
}

type ExpandItemHandlerResponse struct {
}

type ExpandItemHandler struct {
	menuReadModel readmodel.MenuReadModel
}

func NewExpandItemHandler(menuReadModel readmodel.MenuReadModel) *ExpandItemHandler {
	return &ExpandItemHandler{menuReadModel: menuReadModel}
}
func (h *ExpandItemHandler) Handle(ctx context.Context, cmd *ExpandItemCommand) (*ExpandItemHandlerResponse, error) {
	err := h.menuReadModel.ExpandItem(cmd.UserId, cmd.ProjectId, cmd.MenuItemId)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
