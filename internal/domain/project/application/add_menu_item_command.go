package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
	"github.com/google/uuid"
)

type AddMenuItemCommand struct {
	UserId    string
	ProjectId string
	ParentId  string
	NewItemId string
}

type AddMenuItemHandlerResponse struct {
	CanAccess bool
}

type AddMenuItemHandler struct {
	menuRepository      domain.MenuRepository
	editorRepository    domain.EditorRepository
	locationRepository  domain.LocationRepository
	characterRepository domain.CharacterRepository
}

func NewAddMenuItemHandler(
	menuRepository domain.MenuRepository,
	editorRepository domain.EditorRepository,
	locationRepository domain.LocationRepository,
	characterRepository domain.CharacterRepository,
) *AddMenuItemHandler {
	return &AddMenuItemHandler{
		menuRepository:      menuRepository,
		editorRepository:    editorRepository,
		locationRepository:  locationRepository,
		characterRepository: characterRepository,
	}
}

func (h *AddMenuItemHandler) Handle(ctx context.Context, cmd *AddMenuItemCommand) (*AddMenuItemHandlerResponse, error) {
	projectUUid := domain.ProjectId(uuid.MustParse(cmd.ProjectId))
	newMenuItemUUid := domain.MenuItemId(uuid.MustParse(cmd.NewItemId))
	parentUUid := domain.MenuItemId(uuid.MustParse(cmd.ParentId))
	userUUid := domain.UserId(uuid.MustParse(cmd.UserId))

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

	newMenuItem, err := menu.AddNewItem(newMenuItemUUid, parentUUid)

	if err != nil {
		return nil, err
	}

	err = h.menuRepository.Save(tx, userUUid, projectUUid, menu)

	if err != nil {
		return nil, err
	}

	if newMenuItem.Type == domain.TypeChapter {
		editor := domain.NewEditor(newMenuItemUUid)

		err = h.editorRepository.AddEditor(tx, editor)

		if err != nil {
			return nil, err
		}
	}

	if newMenuItem.Type == domain.TypeCharacter {
		character := domain.NewCharacter(newMenuItemUUid)

		err = h.characterRepository.AddCharacter(tx, character)

		if err != nil {
			return nil, err
		}
	}

	if newMenuItem.Type == domain.TypeLocation {
		location := domain.NewLocation(newMenuItemUUid)

		err = h.locationRepository.AddLocation(tx, location)

		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
