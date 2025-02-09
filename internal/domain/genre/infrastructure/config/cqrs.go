package config

import (
	"github.com/gbooksoft/backend/internal/domain/genre/application"
	"github.com/gbooksoft/backend/internal/domain/genre/infrastructure"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/mehdihadeli/go-mediatr"
)

func RegisterMediatr() {
	dbPool := database.GetDBPool()

	genreReadModel := infrastructure.NewMySqlGenreReadModel(dbPool)

	getAllGenresHandler := application.NewGetAllGenresHandler(genreReadModel)

	hE(mediatr.RegisterRequestHandler[*application.GetAllGenresQuery, *application.GetAllGenresHandlerResponse](getAllGenresHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
