package domain

type ProjectService interface {
	Publish(authorId AuthorId, projectId ProjectId) error
}
