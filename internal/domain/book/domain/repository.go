package domain

import (
	"context"
	"github.com/gbooksoft/backend/pkg/repository"
)

type BookRepository interface {
	repository.BaseRepository
	Add(ctx context.Context, book *Book) error
	Save(ctx context.Context, book *Book) error
	FindById(ctx context.Context, id BookId) (*Book, error)
}
