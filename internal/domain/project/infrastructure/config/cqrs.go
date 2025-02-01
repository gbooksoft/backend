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
	getAllProjectsHandler := application.NewGetAllProjectsHandler(projectReadModel)
	getRecentProjectsHandler := application.NewGetRecentProjectsQueryHandler(projectReadModel)
	canAccessHandler := application.NewCanAccessHandler(projectRepository)
	getMetadataHandler := application.NewGetMetadataHandler(projectReadModel)
	updateMetadataHandler := application.NewUpdateMetadataHandler(projectRepository)
	getMenuHandler := application.NewGetMenuHandler(menuReadModel)
	getExpandedItemsHandler := application.NewGetExpandedItemsHandler(menuReadModel)

	hE(mediatr.RegisterRequestHandler[*application.CreateProjectCommand, *application.CreateProjectHandlerResponse](createProjectHandler))
	hE(mediatr.RegisterRequestHandler[*application.ExpandItemCommand, *application.ExpandItemHandlerResponse](expandItemHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetAllProjectsQuery, *application.GetAllProjectsHandlerResponse](getAllProjectsHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetRecentProjectsQuery, *application.GetRecentProjectsHandlerResponse](getRecentProjectsHandler))
	hE(mediatr.RegisterRequestHandler[*application.CanAccessCommand, *application.CanAccessHandlerResponse](canAccessHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetMetadataQuery, *application.GetMetadataHandlerResponse](getMetadataHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateMetadataCommand, *application.UpdateMetadataHandlerResponse](updateMetadataHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetMenuQuery, *application.GetMenuHandlerResponse](getMenuHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetExpandedItemsQuery, *application.GetExpandedItemsHandlerResponse](getExpandedItemsHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
