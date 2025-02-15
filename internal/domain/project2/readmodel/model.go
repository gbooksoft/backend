package readmodel

type Chapter struct {
	Id      string
	Name    string
	Content string
	Order   int
}
type Genre struct {
	Id   string
	Name string
}
type Character struct {
	Id          string
	Name        string
	Description *string
	Attributes  map[string]string
}

type Metadata struct {
	Author   string
	Keywords string
	Language string
	Isbn     string
}
type Location struct {
	Id          string
	Name        string
	Description *string
}
type Project struct {
	Id          string
	Title       string
	Description string
	Genres      []Genre
	Locations   []Location
	Chapters    []Chapter
	Characters  []Character
	Metadata    Metadata
}

type ProjectReadModel interface {
	GetAllLoose(authorId string) ([]Project, error)
	GetById(authorId string, projectId string) (Project, error)
}
