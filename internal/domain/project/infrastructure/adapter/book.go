package adapter

import (
	"context"
	bookDomain "github.com/gbooksoft/backend/internal/domain/book/domain"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type BookAdapter struct {
	bookRepository bookDomain.BookRepository
}

func NewBookAdapter(bookRepository bookDomain.BookRepository) *BookAdapter {
	return &BookAdapter{
		bookRepository: bookRepository,
	}
}

func (a *BookAdapter) Publish(ctx context.Context, project *domain.Project, chaptersToPublish []domain.ChapterInternalId) (*domain.BookId, error) {
	bookId := bookDomain.BookId(uuid.New())

	var genres []bookDomain.Genre
	for _, genre := range project.Genres {
		genres = append(genres, bookDomain.Genre{
			Id:   bookDomain.GenreId(uuid.MustParse(genre.Id.String())),
			Name: genre.Name,
		})
	}

	var chapters []bookDomain.Chapter
	for i, chapter := range project.Chapters {
		if !chapter.IsLastVersion || chapter.DeletedAt != nil {
			continue
		}

		if !a.shouldPublishChapter(chapter, chaptersToPublish) {
			continue
		}

		bookChapterId := bookDomain.ChapterId(uuid.MustParse(chapter.Id.String()))
		chapters = append(chapters, bookDomain.Chapter{
			Id:      bookChapterId,
			BookId:  bookId,
			Title:   chapter.Title,
			Content: chapter.ContentHtml,
			Order:   chapter.Order,
			Version: chapter.Version,
		})

		projectBookChapterId := domain.BookChapterId(uuid.MustParse(bookChapterId.String()))
		project.Chapters[i].BookChapterId = &projectBookChapterId
	}

	book := bookDomain.NewBook(
		bookId,
		project.Title,
		bookDomain.AuthorId(uuid.MustParse(project.AuthorId.String())),
		genres,
		chapters,
		true,
	)

	err := a.bookRepository.Add(ctx, book)

	if err != nil {
		return nil, err
	}

	projectBookId := domain.BookId(uuid.MustParse(bookId.String()))

	return &projectBookId, nil
}

func (a *BookAdapter) Republish(ctx context.Context, project *domain.Project, chaptersToPublish []domain.ChapterInternalId) error {
	bookId := bookDomain.BookId(uuid.MustParse(project.BookId.String()))
	book, err := a.bookRepository.FindById(ctx, bookId)
	if err != nil {
		return err
	}

	if book == nil {
		return gerr.NewNotFoundErr("published book not found")
	}

	var renewedChapters []bookDomain.Chapter

	for i, chapter := range project.Chapters {
		if !chapter.IsLastVersion {
			continue
		}

		if chapter.DeletedAt != nil {
			continue
		}

		if !a.shouldPublishChapter(chapter, chaptersToPublish) {
			continue
		}

		var bookChapter bookDomain.Chapter
		bookChapter, ok := a.findChapter(book, chapter)

		if !ok {
			bookChapterId := bookDomain.ChapterId(uuid.MustParse(chapter.Id.String()))
			bookChapter = bookDomain.Chapter{
				Id:      bookChapterId,
				BookId:  bookId,
				Title:   chapter.Title,
				Content: chapter.ContentHtml,
				Order:   chapter.Order,
				Version: chapter.Version,
			}
			projectBookChapterId := domain.BookChapterId(uuid.MustParse(bookChapterId.String()))
			project.Chapters[i].BookChapterId = &projectBookChapterId
		} else {
			bookChapter.Title = chapter.Title
			bookChapter.Content = chapter.ContentHtml
			bookChapter.Order = chapter.Order
			bookChapter.Version = chapter.Version
		}

		renewedChapters = append(renewedChapters, bookChapter)
	}

	book.Chapters = renewedChapters

	a.unTouchChapterRelation(project, renewedChapters)

	err = a.bookRepository.Save(ctx, book)

	if err != nil {
		return err
	}

	return nil
}

func (a *BookAdapter) shouldPublishChapter(chapter domain.Chapter, chaptersToPublish []domain.ChapterInternalId) bool {
	if len(chaptersToPublish) == 0 {
		return true
	}

	for _, c := range chaptersToPublish {
		if c.Equal(chapter.ChapterInternalId) {
			return true
		}
	}

	return false
}

func (a *BookAdapter) findChapter(book *bookDomain.Book, chapter domain.Chapter) (bookDomain.Chapter, bool) {
	if chapter.BookChapterId == nil {
		return bookDomain.Chapter{}, false
	}

	for _, ch := range book.Chapters {
		if ch.Id.String() == chapter.BookChapterId.String() {
			r := ch
			return r, true
		}
	}

	return bookDomain.Chapter{}, false
}

func (a *BookAdapter) unTouchChapterRelation(project *domain.Project, chapters []bookDomain.Chapter) {
	for i, c := range project.Chapters {
		if c.BookChapterId == nil {
			continue
		}

		found := false

		for _, ch := range chapters {
			if c.BookChapterId.String() == ch.Id.String() {
				found = true
				break
			}
		}

		if !found {
			project.Chapters[i].BookChapterId = nil
		}
	}
}
