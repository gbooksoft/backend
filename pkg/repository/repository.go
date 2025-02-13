package repository

import (
	"context"
	"github.com/google/uuid"
)

type BaseRepository interface {
	GetNextId() uuid.UUID
	BeginTransaction(ctx context.Context) (context.Context, error)
	RollbackTransaction(ctx context.Context) error
	CommitTransaction(ctx context.Context) error
}
