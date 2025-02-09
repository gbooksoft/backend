package domain

import "github.com/google/uuid"

type UserId uuid.UUID

func (id UserId) String() string {
	return uuid.UUID(id).String()
}

type ProjectId uuid.UUID

func (id ProjectId) String() string {
	return uuid.UUID(id).String()
}

type Project struct {
	Id          ProjectId
	UserId      UserId
	Title       string
	Description string
	GenresIds   []GenreId
	Tags        []string
}

func NewProject(id ProjectId, userId UserId, title, description string, GenresIds []GenreId, tags []string) *Project {
	return &Project{
		Id:          id,
		UserId:      userId,
		Title:       title,
		Description: description,
		GenresIds:   GenresIds,
		Tags:        tags,
	}
}

type MetadataId uuid.UUID

func (id MetadataId) String() string {
	return uuid.UUID(id).String()
}

type Metadata struct {
	Id        MetadataId
	ProjectId ProjectId
	Author    string
	Genre     string
	Keywords  string
	Language  string
	Isbn      string
}

func NewMetadata(id MetadataId, projectId ProjectId, author string, genre string, keywords string, language string, isbn string) *Metadata {
	return &Metadata{
		Id:        id,
		ProjectId: projectId,
		Author:    author,
		Genre:     genre,
		Keywords:  keywords,
		Language:  language,
		Isbn:      isbn,
	}
}
