package readmodel

type Project struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ProjectReadModel interface {
	GetAll(userId string) ([]Project, error)
	GetRecent(userId string, limit int) ([]Project, error)
	GetById(userId string, projectId string) (Project, error)
	GetMetadata(userId string, projectId string) (*Metadata, error)
}
