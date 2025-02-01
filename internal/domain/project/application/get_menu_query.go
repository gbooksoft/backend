package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
	"github.com/gbooksoft/backend/pkg/gerr"
)

type GetMenuQuery struct {
	UserId    string
	ProjectId string
}

type GetMenuHandlerResponse struct {
	Menu *readmodel.Menu
}
type GetMenuHandler struct {
	menuReadModel readmodel.MenuReadModel
}

func NewGetMenuHandler(menuReadModel readmodel.MenuReadModel) *GetMenuHandler {
	return &GetMenuHandler{menuReadModel: menuReadModel}
}

func (h *GetMenuHandler) Handle(ctx context.Context, query *GetMenuQuery) (*GetMenuHandlerResponse, error) {
	menu, err := h.menuReadModel.GetByUserIdAndProjectId(query.UserId, query.ProjectId)

	if err != nil {
		return nil, err
	}

	if menu == nil || len(menu.Items) == 0 {
		return nil, gerr.NewBusinessLogicErr("Menu not found", gerr.NotFound)
	}

	return &GetMenuHandlerResponse{Menu: menu}, nil
}
