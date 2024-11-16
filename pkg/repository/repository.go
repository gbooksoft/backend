package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type BaseRepository interface {
	GetNextId() uuid.UUID
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
}

type BaseMysqlRepository struct {
	db *sql.DB
}

func NewBaseMysqlRepository(db *sql.DB) *BaseMysqlRepository {
	return &BaseMysqlRepository{
		db: db,
	}
}

func (b *BaseMysqlRepository) GetNextId() uuid.UUID {
	return uuid.New()
}
func (b *BaseMysqlRepository) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return b.db.BeginTx(ctx, nil)
}
