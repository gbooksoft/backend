package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/google/uuid"
)

type DeleteMenuItemCommand struct {
	UserId     string
	ProjectId  string
	MenuItemId string
}

type DeleteMenuItemHandlerResponse struct {
}

type DeleteMenuItemHandler struct {
	menuRepository domain.MenuRepository
}

func NewDeleteMenuItemHandler(menuRepository domain.MenuRepository) *DeleteMenuItemHandler {
	return &DeleteMenuItemHandler{
		menuRepository: menuRepository,
	}
}

func (h *DeleteMenuItemHandler) Handle(ctx context.Context, cmd *DeleteMenuItemCommand) (*DeleteMenuItemHandlerResponse, error) {
	userUUid := domain.UserId(uuid.MustParse(cmd.UserId))
	projectUUid := domain.ProjectId(uuid.MustParse(cmd.ProjectId))
	menuItemUuid := domain.MenuItemId(uuid.MustParse(cmd.MenuItemId))

	err := h.menuRepository.DeleteMenuItem(userUUid, projectUUid, menuItemUuid)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
