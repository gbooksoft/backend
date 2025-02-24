package infrastructure

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/book/domain"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/gbooksoft/backend/pkg/repository"
)

type MySqlBookRepository struct {
	*repository.BaseGormRepository
}

func NewMySqlBookRepository() *MySqlBookRepository {
	return &MySqlBookRepository{
		repository.NewBaseGormRepository(database.GetGormDB()),
	}
}

func (r *MySqlBookRepository) Add(ctx context.Context, book *domain.Book) error {
	tx, ok := repository.FetchGormTransactionFromContext(ctx)

	db := r.Db
	if ok {
		db = tx
	}

	result := db.Table("books").Create(book)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}
