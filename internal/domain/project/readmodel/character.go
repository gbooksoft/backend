package readmodel

type Attribute struct {
	Field string `json:"field"`
	Value string `json:"value"`
}
type Character struct {
	MenuItemId  string
	Name        string
	Description string
	Attributes  []Attribute `json:"attributes"`
}

type CharacterReadModel interface {
	GetCharacters(userId, projectId string) ([]Character, error)
}
