package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
	"github.com/google/uuid"
)

type UpdateMenuItemCommand struct {
	UserId      string
	ProjectId   string
	MenuItemId  string
	Title       string
	SetSelected bool
}

type UpdateMenuItemHandlerResponse struct {
}

type UpdateMenuItemHandler struct {
	menuRepository domain.MenuRepository
}

func NewUpdateMenuItemHandler(menuRepository domain.MenuRepository) *UpdateMenuItemHandler {
	return &UpdateMenuItemHandler{
		menuRepository: menuRepository,
	}
}

func (h *UpdateMenuItemHandler) Handle(ctx context.Context, cmd *UpdateMenuItemCommand) (*UpdateMenuItemHandlerResponse, error) {
	userUUid := domain.UserId(uuid.MustParse(cmd.UserId))
	projectUUid := domain.ProjectId(uuid.MustParse(cmd.ProjectId))
	menuItemUuid := domain.MenuItemId(uuid.MustParse(cmd.MenuItemId))
	title := cmd.Title
	setSelected := cmd.SetSelected

	ctx, err := h.menuRepository.BeginTransaction(context.TODO())
	if err != nil {
		return nil, err
	}
	tx, ok := repository.FetchDBTransactionFromContext(ctx)
	if !ok {
		panic("transaction not found")
	}
	defer tx.Rollback()

	menu, err := h.menuRepository.FindByProjectId(tx, userUUid, projectUUid)
	if err != nil {
		return nil, err
	}

	if len(title) > 0 {
		err = menu.UpdateItemTitle(menuItemUuid, title)

		if err != nil {
			return nil, err
		}
	}

	if setSelected {
		menu.SelectItem(menuItemUuid)
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
