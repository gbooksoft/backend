package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/genre/readmodel"
)

type GetAllGenresQuery struct {
	Name string
}

type GetAllGenresHandlerResponse struct {
	Genres []readmodel.Genre
}

type GetAllGenresHandler struct {
	genreReadModel readmodel.GenreReadModel
}

func NewGetAllGenresHandler(genreReadModel readmodel.GenreReadModel) *GetAllGenresHandler {
	return &GetAllGenresHandler{genreReadModel: genreReadModel}
}

func (h *GetAllGenresHandler) Handle(ctx context.Context, query *GetAllGenresQuery) (*GetAllGenresHandlerResponse, error) {
	genres, err := h.genreReadModel.GetAllByName(query.Name)
	if err != nil {
		return nil, err
	}

	return &GetAllGenresHandlerResponse{Genres: genres}, nil
}
