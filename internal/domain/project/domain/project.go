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
	Id     ProjectId
	UserId UserId
	Title  string
}

func NewProject(id ProjectId, userId UserId, title string) *Project {
	return &Project{
		Id:     id,
		UserId: userId,
		Title:  title,
	}
}
