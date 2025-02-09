package readmodel

type Genre struct {
	Id   string
	Name string
}

type GenreReadModel interface {
	GetAllByName(name string) ([]Genre, error)
}
