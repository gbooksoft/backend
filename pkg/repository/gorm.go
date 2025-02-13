package repository

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormTxKey struct{}

type BaseGormRepository struct {
	Db *gorm.DB
}

func NewBaseGormRepository(db *gorm.DB) *BaseGormRepository {
	return &BaseGormRepository{
		Db: db,
	}
}

func (b *BaseGormRepository) GetNextId() uuid.UUID {
	return uuid.New()
}

func (b *BaseGormRepository) BeginTransaction(ctx context.Context) (context.Context, error) {
	txDb := b.Db.Begin()

	if err := txDb.Error; err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, gormTxKey{}, txDb), nil
}

func (b *BaseGormRepository) RollbackTransaction(ctx context.Context) error {
	tx, ok := FetchGormTransactionFromContext(ctx)

	if !ok {
		return nil
	}

	if err := tx.Rollback().Error; err != nil {
		return err
	}

	return nil
}

func (b *BaseGormRepository) CommitTransaction(ctx context.Context) error {
	tx, ok := FetchGormTransactionFromContext(ctx)

	if !ok {
		return nil
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func FetchGormTransactionFromContext(ctx context.Context) (*gorm.DB, bool) {
	e, ok := ctx.Value(gormTxKey{}).(*gorm.DB)
	if !ok {
		return nil, false
	}

	return e, true
}
