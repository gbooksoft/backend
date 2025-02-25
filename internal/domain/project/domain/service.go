package domain

import "context"

type ProjectService interface {
	Publish(ctx context.Context, authorId AuthorId, projectId ProjectId, chapters []ChapterInternalId) error
	Republish(ctx context.Context, authorId AuthorId, projectId ProjectId, chapters []ChapterInternalId) error
}
