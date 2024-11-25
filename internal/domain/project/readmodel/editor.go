package readmodel

type Editor struct {
	MenuItemId  string
	EditorState string
}

type EditorReadModel interface {
	GetEditors(userId, projectId string) ([]Editor, error)
}
