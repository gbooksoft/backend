package adapter

import (
	"context"
	"database/sql"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
)

type GenreAdapter struct {
	*repository.BaseMysqlRepository
}

func NewGenreAdapter(db *sql.DB) *GenreAdapter {
	return &GenreAdapter{repository.NewBaseMysqlRepository(db)}
}

func (g *GenreAdapter) Exists(ctx context.Context, genres []domain.GenreId) (bool, error) {
	return false, nil
}
