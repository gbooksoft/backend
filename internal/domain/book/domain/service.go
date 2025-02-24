package domain

import "context"

type PublishService interface {
	PublishBook(ctx context.Context, book *Book) error
}
