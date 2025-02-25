package infrastructure

import (
	"context"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/book/domain"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/gbooksoft/backend/pkg/repository"
	"gorm.io/gorm"
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

func (r *MySqlBookRepository) Save(ctx context.Context, book *domain.Book) error {
	tx, ok := repository.FetchGormTransactionFromContext(ctx)

	db := r.Db
	if ok {
		db = tx
	}

	if len(book.Chapters) == 0 {
		if err := db.Unscoped().Model(book).Association("Chapters").Unscoped().Clear(); err != nil {
			return err
		}
	} else {
		var newChapterIds []string
		for _, chapter := range book.Chapters {
			newChapterIds = append(newChapterIds, chapter.Id.String())
		}
		if err := db.Where("book_id = ? AND id NOT IN ?", book.Id.String(), newChapterIds).Delete(&domain.Chapter{}).Error; err != nil {
			return err
		}
	}

	result := db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(book)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func (r *MySqlBookRepository) FindById(ctx context.Context, id domain.BookId) (*domain.Book, error) {
	tx, ok := repository.FetchGormTransactionFromContext(ctx)

	db := r.Db
	if ok {
		db = tx
	}

	book := &domain.Book{}

	if err := db.
		Preload("Genres").
		Preload("Chapters").
		First(book, "id = ?", id.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return book, nil
}
