package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/google/uuid"
)

type SaveEditorStateCommand struct {
	UserId      string
	ProjectId   string
	MenuItemId  string
	EditorState string
}

type SaveEditorStateHandlerResponse struct {
}

type SaveEditorStateHandler struct {
	editorRepository domain.EditorRepository
}

func NewSaveEditorStateHandler(editorRepository domain.EditorRepository) *SaveEditorStateHandler {
	return &SaveEditorStateHandler{
		editorRepository: editorRepository,
	}
}

func (h *SaveEditorStateHandler) Handle(ctx context.Context, cmd *SaveEditorStateCommand) (*SaveEditorStateHandlerResponse, error) {
	userUUid := domain.UserId(uuid.MustParse(cmd.UserId))
	projectUUid := domain.ProjectId(uuid.MustParse(cmd.ProjectId))
	menuItemUuid := domain.MenuItemId(uuid.MustParse(cmd.MenuItemId))
	editorState := cmd.EditorState

	tx, err := h.editorRepository.BeginTransaction(context.TODO())

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	editor, err := h.editorRepository.FindEditor(tx, userUUid, projectUUid, menuItemUuid)

	if err != nil {
		return nil, err
	}

	editor.EditorState = &editorState

	err = h.editorRepository.Save(tx, editor)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
