package readmodel

type Menu struct {
	Items []MenuItem
}

type MenuItem struct {
	Id         string
	Title      string
	ItemType   string
	ParentId   *string
	IsSelected bool
	Order      int
}

type MenuReadModel interface {
	GetByUserIdAndProjectId(userId, projectId string) (*Menu, error)
	GetExpandedItems(userId string, projectId string) ([]string, error)
	ExpandItem(userId string, projectId string, menuItemId string) error
	CollapseItem(userId string, projectId string, menuItemId string) error
}
