package domain

import (
	"context"
	"github.com/gbooksoft/backend/pkg/repository"
)

type BookRepository interface {
	repository.BaseRepository
	Add(ctx context.Context, book *Book) error
}
