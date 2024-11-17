package readmodel

type Menu struct {
	Items []MenuItem
}

type MenuItem struct {
	Id       string
	Title    string
	ItemType string
	ParentId *string
}

type MenuReadModel interface {
	GetByUserIdAndProjectId(userId, projectId string) (*Menu, error)
	GetExpandedItems(userId string, projectId string) []string
}
