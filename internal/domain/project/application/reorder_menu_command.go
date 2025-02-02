package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/google/uuid"
)

type MenuItem struct {
	ItemId   string
	ParentId string
	Order    int
}

type ReorderMenuCommand struct {
	UserId    string
	ProjectId string
	MenuItems []MenuItem
}

type ReorderMenuHandlerResponse struct {
}

type ReorderMenuHandler struct {
	menuRepository domain.MenuRepository
}

func NewReorderMenuHandler(menuRepository domain.MenuRepository) *ReorderMenuHandler {
	return &ReorderMenuHandler{
		menuRepository: menuRepository,
	}
}

func (h *ReorderMenuHandler) Handle(ctx context.Context, cmd *ReorderMenuCommand) (*ReorderMenuHandlerResponse, error) {
	userUUid := domain.UserId(uuid.MustParse(cmd.UserId))
	projectUUid := domain.ProjectId(uuid.MustParse(cmd.ProjectId))

	tx, err := h.menuRepository.BeginTransaction(context.TODO())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	menu, err := h.menuRepository.FindByProjectId(tx, userUUid, projectUUid)
	if err != nil {
		return nil, err
	}

	var itemsToReorder []domain.ReorderItem
	for _, item := range cmd.MenuItems {
		itemsToReorder = append(itemsToReorder, domain.ReorderItem{
			ItemId:       item.ItemId,
			ParentItemId: item.ParentId,
			Order:        item.Order,
		})
	}

	err = menu.Reorder(itemsToReorder)

	if err != nil {
		return nil, err
	}

	err = h.menuRepository.Save(tx, userUUid, projectUUid, menu)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
