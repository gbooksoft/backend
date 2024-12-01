package domain

type Editor struct {
	MenuItemId  MenuItemId
	EditorState *string
}

func NewEditor(menuItemId MenuItemId) *Editor {
	return &Editor{
		MenuItemId: menuItemId,
	}
}

func NewEditorsFromMenu(menu *Menu) []*Editor {
	items := menu.findChildrenWithType(TypeChapter)

	editors := make([]*Editor, len(items))
	for index, i := range items {
		editor := NewEditor(i.Id)

		editors[index] = editor
	}

	return editors
}
