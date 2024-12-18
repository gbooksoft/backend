package domain

type Location struct {
	MenuItemId  MenuItemId
	Name        string
	Description *string
}

func NewLocation(menuItemId MenuItemId) *Location {
	return &Location{
		MenuItemId: menuItemId,
		Name:       "Название локации",
	}
}
