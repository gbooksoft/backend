package readmodel

type Location struct {
	MenuItemId  string
	Name        string
	Description string
}

type LocationReadModel interface {
	GetLocations(userId, projectId string) ([]Location, error)
}
