package config

import (
	"github.com/gbooksoft/backend/internal/domain/project2/application"
	"github.com/gbooksoft/backend/internal/domain/project2/infrastructure"
	"github.com/gbooksoft/backend/internal/domain/project2/infrastructure/adapter"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/mehdihadeli/go-mediatr"
)

func RegisterMediatr() {
	dbPool := database.GetDBPool()

	projectRepository := infrastructure.NewMySqlProjectRepository(dbPool)
	genreAdapter := adapter.NewGenreAdapter(dbPool)
	createProjectHandler := application.NewCreateProjectHandler(projectRepository, genreAdapter)

	hE(mediatr.RegisterRequestHandler[*application.CreateProjectCommand, *application.CreateProjectHandlerResponse](createProjectHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
