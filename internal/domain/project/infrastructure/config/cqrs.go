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
	locationRepository := infrastructure.NewMySqlLocationRepository(dbPool)
	characterRepository := infrastructure.NewMySqlCharacterRepository(dbPool)
	editorReadModel := infrastructure.NewMySqlEditorReadModel(dbPool)
	characterReadModel := infrastructure.NewMySqlCharacterReadModel(dbPool)
	locationReadModel := infrastructure.NewMySqlLocationReadModel(dbPool)

	createProjectHandler := application.NewCreateProjectHandler(projectRepository, menuRepository, editorRepository)
	expandItemHandler := application.NewExpandItemHandler(menuReadModel)
	getAllProjectsHandler := application.NewGetAllProjectsHandler(projectReadModel)
	getRecentProjectsHandler := application.NewGetRecentProjectsQueryHandler(projectReadModel)
	canAccessHandler := application.NewCanAccessHandler(projectRepository)
	getMetadataHandler := application.NewGetMetadataHandler(projectReadModel)
	updateMetadataHandler := application.NewUpdateMetadataHandler(projectRepository)
	getMenuHandler := application.NewGetMenuHandler(menuReadModel)
	getExpandedItemsHandler := application.NewGetExpandedItemsHandler(menuReadModel)
	addMenuItemHandler := application.NewAddMenuItemHandler(menuRepository, editorRepository, locationRepository, characterRepository)
	collapseItemHandler := application.NewCollapseItemHandler(menuReadModel)
	updateMenuItemHandler := application.NewUpdateMenuItemHandler(menuRepository)
	deleteMenuItemHandler := application.NewDeleteMenuItemHandler(menuRepository)
	getEditorsHandler := application.NewGetEditorsHandler(editorReadModel)
	saveEditorStateHandler := application.NewSaveEditorStateHandler(editorRepository)
	getCharactersHandler := application.NewGetCharactersHandler(characterReadModel)
	updateCharacterHandler := application.NewUpdateCharacterHandler(characterRepository)
	getLocationsHandler := application.NewGetLocationsHandler(locationReadModel)
	updateLocationHandler := application.NewUpdateLocationHandler(locationRepository)
	getProjectHandler := application.NewGetProjectHandler(projectReadModel)
	updateProjectHandler := application.NewUpdateProjectHandler(projectRepository)

	hE(mediatr.RegisterRequestHandler[*application.CreateProjectCommand, *application.CreateProjectHandlerResponse](createProjectHandler))
	hE(mediatr.RegisterRequestHandler[*application.ExpandItemCommand, *application.ExpandItemHandlerResponse](expandItemHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetAllProjectsQuery, *application.GetAllProjectsHandlerResponse](getAllProjectsHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetRecentProjectsQuery, *application.GetRecentProjectsHandlerResponse](getRecentProjectsHandler))
	hE(mediatr.RegisterRequestHandler[*application.CanAccessCommand, *application.CanAccessHandlerResponse](canAccessHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetMetadataQuery, *application.GetMetadataHandlerResponse](getMetadataHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateMetadataCommand, *application.UpdateMetadataHandlerResponse](updateMetadataHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetMenuQuery, *application.GetMenuHandlerResponse](getMenuHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetExpandedItemsQuery, *application.GetExpandedItemsHandlerResponse](getExpandedItemsHandler))
	hE(mediatr.RegisterRequestHandler[*application.AddMenuItemCommand, *application.AddMenuItemHandlerResponse](addMenuItemHandler))
	hE(mediatr.RegisterRequestHandler[*application.CollapseItemCommand, *application.CollapseItemHandlerResponse](collapseItemHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateMenuItemCommand, *application.UpdateMenuItemHandlerResponse](updateMenuItemHandler))
	hE(mediatr.RegisterRequestHandler[*application.DeleteMenuItemCommand, *application.DeleteMenuItemHandlerResponse](deleteMenuItemHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetEditorsQuery, *application.GetEditorsHandlerResponse](getEditorsHandler))
	hE(mediatr.RegisterRequestHandler[*application.SaveEditorStateCommand, *application.SaveEditorStateHandlerResponse](saveEditorStateHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetCharactersQuery, *application.GetCharactersHandlerResponse](getCharactersHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateCharacterCommand, *application.UpdateCharacterHandlerResponse](updateCharacterHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetLocationsQuery, *application.GetLocationsHandlerResponse](getLocationsHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateLocationCommand, *application.UpdateLocationHandlerResponse](updateLocationHandler))
	hE(mediatr.RegisterRequestHandler[*application.GetProjectQuery, *application.GetProjectHandlerResponse](getProjectHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateProjectCommand, *application.UpdateProjectHandlerResponse](updateProjectHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
