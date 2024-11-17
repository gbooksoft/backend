package domain

import "github.com/google/uuid"

type EditorId uuid.UUID

func (id EditorId) String() string {
	return uuid.UUID(id).String()
}

type Editor struct {
	Id          EditorId
	MenuItemId  MenuItemId
	EditorState *string
}

func NewEditor(menuItemId MenuItemId) *Editor {
	return &Editor{
		Id:         EditorId(uuid.New()),
		MenuItemId: menuItemId,
	}
}

func NewEditorsFromMenu(menu *Menu) []*Editor {
	items := menu.findChildrenWithType(Chapter)

	editors := make([]*Editor, len(items))
	for index, i := range items {
		editor := NewEditor(i.Id)

		editors[index] = editor
	}

	return editors
}
