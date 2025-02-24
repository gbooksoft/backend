package config

import (
	bookInfrastructure "github.com/gbooksoft/backend/internal/domain/book/infrastructure"
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/internal/domain/project/infrastructure"
	"github.com/gbooksoft/backend/internal/domain/project/infrastructure/adapter"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/mehdihadeli/go-mediatr"
)

func RegisterMediatr() {
	gormDb := database.GetGormDB()

	projectReadModel := infrastructure.NewMySqlProjectReadModel(gormDb)
	projectRepository := infrastructure.NewMySqlProjectRepository()
	bookRepository := bookInfrastructure.NewMySqlBookRepository()
	bookAdapter := adapter.NewBookAdapter(bookRepository)

	projectService := infrastructure.NewProjectService(projectRepository, bookAdapter)

	createProjectHandler := application.NewCreateProjectHandler(projectRepository)
	updateProjectHandler := application.NewUpdateProjectHandler(projectRepository)

	createChapterHandler := application.NewCreateChapterHandler(projectRepository)
	updateChapterHandler := application.NewUpdateChapterHandler(projectRepository)
	deleteChapterHandler := application.NewDeleteChapterHandler(projectRepository)
	reorderChapterHandler := application.NewReorderChapterHandler(projectRepository)

	createCharacterHandler := application.NewCreateCharacterHandler(projectRepository)
	updateCharacterHandler := application.NewUpdateCharacterHandler(projectRepository)
	deleteCharacterHandler := application.NewDeleteCharacterHandler(projectRepository)

	createLocationHandler := application.NewCreateLocationHandler(projectRepository)
	updateLocationHandler := application.NewUpdateLocationHandler(projectRepository)
	deleteLocationHandler := application.NewDeleteLocationHandler(projectRepository)

	updateMetadataHandler := application.NewUpdateMetadataHandler(projectRepository)

	getAllProjectsQuery := application.NewGetAllProjectsHandler(projectReadModel)
	getProjectsQuery := application.NewGetProjectHandler(projectReadModel)

	publishProjectHandler := application.NewPublishProjectHandler(projectService)
	republishProjectHandler := application.NewRepublishProjectHandler(projectService)

	hE(mediatr.RegisterRequestHandler[*application.CreateProjectCommand, *application.CreateProjectHandlerResponse](createProjectHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateProjectCommand, *application.UpdateProjectHandlerResponse](updateProjectHandler))

	hE(mediatr.RegisterRequestHandler[*application.CreateChapterCommand, *application.CreateChapterHandlerResponse](createChapterHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateChapterCommand, *application.UpdateChapterHandlerResponse](updateChapterHandler))
	hE(mediatr.RegisterRequestHandler[*application.DeleteChapterCommand, *application.DeleteChapterHandlerResponse](deleteChapterHandler))
	hE(mediatr.RegisterRequestHandler[*application.ReorderChapterCommand, *application.ReorderChapterHandlerResponse](reorderChapterHandler))

	hE(mediatr.RegisterRequestHandler[*application.CreateCharacterCommand, *application.CreateCharacterHandlerResponse](createCharacterHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateCharacterCommand, *application.UpdateCharacterHandlerResponse](updateCharacterHandler))
	hE(mediatr.RegisterRequestHandler[*application.DeleteCharacterCommand, *application.DeleteCharacterHandlerResponse](deleteCharacterHandler))

	hE(mediatr.RegisterRequestHandler[*application.CreateLocationCommand, *application.CreateLocationHandlerResponse](createLocationHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateLocationCommand, *application.UpdateLocationHandlerResponse](updateLocationHandler))
	hE(mediatr.RegisterRequestHandler[*application.DeleteLocationCommand, *application.DeleteLocationHandlerResponse](deleteLocationHandler))

	hE(mediatr.RegisterRequestHandler[*application.UpdateMetadataCommand, *application.UpdateMetadataHandlerResponse](updateMetadataHandler))

	hE(mediatr.RegisterRequestHandler[*application.GetAllProjectsQuery, *application.GetAllProjectsHandlerResponse](getAllProjectsQuery))
	hE(mediatr.RegisterRequestHandler[*application.GetProjectQuery, *application.GetProjectHandlerResponse](getProjectsQuery))

	hE(mediatr.RegisterRequestHandler[*application.PublishProjectCommand, *application.PublishProjectHandlerResponse](publishProjectHandler))
	hE(mediatr.RegisterRequestHandler[*application.RepublishProjectCommand, *application.RepublishProjectHandlerResponse](republishProjectHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
