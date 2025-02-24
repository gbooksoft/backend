package infrastructure

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/book/domain"
)

type PublishService struct {
	bookRepository domain.BookRepository
}

func NewPublishService(bookRepository domain.BookRepository) *PublishService {
	return &PublishService{
		bookRepository: bookRepository,
	}
}

func (s *PublishService) PublishBook(ctx context.Context, book *domain.Book) error {
	return s.bookRepository.Add(ctx, book)
}
