package config

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/internal/domain/project/infrastructure"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/mehdihadeli/go-mediatr"
)

func RegisterMediatr() {
	dbPool := database.GetDBPool()

	projectRepository := infrastructure.NewMySqlProjectRepository(dbPool)
	menuRepository := infrastructure.NewMySqlMenuRepository(dbPool)
	editorRepository := infrastructure.NewMySqlEditorRepository(dbPool)
	menuReadModel := infrastructure.NewMySqlMenuReadModel(dbPool)
	projectReadModel := infrastructure.NewMySqlProjectReadModel(dbPool)

	createProjectHandler := application.NewCreateProjectHandler(projectRepository, menuRepository, editorRepository)
	expandItemHandler := application.NewExpandItemHandler(menuReadModel)
	getAllProjectsHandler := application.NewGetAllProjectsQueryHandler(projectReadModel)

	hE(mediatr.RegisterRequestHandler[*application.CreateProjectCommand, *application.CreateProjectHandlerResponse](createProjectHandler))
	hE(mediatr.RegisterRequestHandler[*application.ExpandItemCommand, *application.ExpandItemHandlerResponse](expandItemHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetAllProjectsQuery, *application.GetAllProjectsQueryHandlerResponse](getAllProjectsHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
