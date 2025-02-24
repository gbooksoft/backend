package adapter

import (
	"context"
	bookDomain "github.com/gbooksoft/backend/internal/domain/book/domain"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/google/uuid"
)

type BookAdapter struct {
	publishService bookDomain.PublishService
}

func NewBookAdapter(publishService bookDomain.PublishService) *BookAdapter {
	return &BookAdapter{
		publishService: publishService,
	}
}

func (a *BookAdapter) Publish(ctx context.Context, project *domain.Project) (*domain.BookId, error) {
	bookId := bookDomain.BookId(uuid.New())

	var genres []bookDomain.Genre
	for _, genre := range project.Genres {
		genres = append(genres, bookDomain.Genre{
			Id:   bookDomain.GenreId(uuid.MustParse(genre.Id.String())),
			Name: genre.Name,
		})
	}

	var chapters []bookDomain.Chapter
	for _, chapter := range project.Chapters {
		if !chapter.IsLastVersion || chapter.DeletedAt != nil {
			continue
		}

		chapters = append(chapters, bookDomain.Chapter{
			Id:      bookDomain.ChapterId(uuid.MustParse(chapter.Id.String())),
			BookId:  bookId,
			Title:   chapter.Title,
			Content: chapter.ContentHtml,
			Order:   chapter.Order,
		})
	}

	book := bookDomain.NewBook(
		bookId,
		project.Title,
		bookDomain.AuthorId(uuid.MustParse(project.AuthorId.String())),
		genres,
		chapters,
		true,
	)

	err := a.publishService.PublishBook(ctx, book)

	if err != nil {
		return nil, err
	}

	projectBookId := domain.BookId(uuid.MustParse(bookId.String()))

	return &projectBookId, nil
}
