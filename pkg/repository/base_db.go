package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

type dbTxKey struct{}

type BaseMysqlRepository struct {
	Db *sql.DB
}

func NewBaseMysqlRepository(db *sql.DB) *BaseMysqlRepository {
	return &BaseMysqlRepository{
		Db: db,
	}
}

func (b *BaseMysqlRepository) GetNextId() uuid.UUID {
	return uuid.New()
}

func (b *BaseMysqlRepository) BeginTransaction(ctx context.Context) (context.Context, error) {
	tx, err := b.Db.BeginTx(ctx, nil)

	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, dbTxKey{}, tx), nil
}

func (b *BaseMysqlRepository) RollbackTransaction(ctx context.Context) error {
	return errors.New("not implemented")
}

func (b *BaseMysqlRepository) CommitTransaction(ctx context.Context) error {
	return errors.New("not implemented")
}

func FetchDBTransactionFromContext(ctx context.Context) (*sql.Tx, bool) {
	e, ok := ctx.Value(dbTxKey{}).(*sql.Tx)
	if !ok {
		return nil, false
	}

	return e, true
}
