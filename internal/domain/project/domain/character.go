package domain

type Attribute struct {
	Field string
	Value string
}

type Character struct {
	MenuItemId  MenuItemId
	Name        string
	Description string
	Attributes  []Attribute
}

func NewCharacter(menuItemId MenuItemId) *Character {
	return &Character{
		MenuItemId:  menuItemId,
		Name:        "Character",
		Description: "Description...",
		Attributes: []Attribute{
			{
				Field: "Возраст",
				Value: "15",
			},
		},
	}
}
