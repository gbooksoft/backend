package readmodel

type Author struct {
	Id         string
	Email      string
	Username   string
	FirstName  string
	LastName   string
	AvatarPath string
}

type AuthorReadModel interface {
	GetByUsername(username string) (Author, error)
}
