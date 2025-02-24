package domain

type ProjectService interface {
	Publish(authorId AuthorId, projectId ProjectId, chapters []ChapterInternalId) error
}
