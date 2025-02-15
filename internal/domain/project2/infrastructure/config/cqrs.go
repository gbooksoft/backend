package config

import (
	"github.com/gbooksoft/backend/internal/domain/project2/application"
	"github.com/gbooksoft/backend/internal/domain/project2/infrastructure"
	"github.com/mehdihadeli/go-mediatr"
)

func RegisterMediatr() {
	projectRepository := infrastructure.NewMySqlProjectRepository()
	createProjectHandler := application.NewCreateProjectHandler(projectRepository)
	updateProjectHandler := application.NewUpdateProjectHandler(projectRepository)

	createChapterHandler := application.NewCreateChapterHandler(projectRepository)
	updateChapterHandler := application.NewUpdateChapterHandler(projectRepository)
	deleteChapterHandler := application.NewDeleteChapterHandler(projectRepository)

	createCharacterHandler := application.NewCreateCharacterHandler(projectRepository)
	updateCharacterHandler := application.NewUpdateCharacterHandler(projectRepository)
	deleteCharacterHandler := application.NewDeleteCharacterHandler(projectRepository)

	createLocationHandler := application.NewCreateLocationHandler(projectRepository)
	updateLocationHandler := application.NewUpdateLocationHandler(projectRepository)
	deleteLocationHandler := application.NewDeleteLocationHandler(projectRepository)

	updateMetadataHandler := application.NewUpdateMetadataHandler(projectRepository)

	hE(mediatr.RegisterRequestHandler[*application.CreateProjectCommand, *application.CreateProjectHandlerResponse](createProjectHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateProjectCommand, *application.UpdateProjectHandlerResponse](updateProjectHandler))

	hE(mediatr.RegisterRequestHandler[*application.CreateChapterCommand, *application.CreateChapterHandlerResponse](createChapterHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateChapterCommand, *application.UpdateChapterHandlerResponse](updateChapterHandler))
	hE(mediatr.RegisterRequestHandler[*application.DeleteChapterCommand, *application.DeleteChapterHandlerResponse](deleteChapterHandler))

	hE(mediatr.RegisterRequestHandler[*application.CreateCharacterCommand, *application.CreateCharacterHandlerResponse](createCharacterHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateCharacterCommand, *application.UpdateCharacterHandlerResponse](updateCharacterHandler))
	hE(mediatr.RegisterRequestHandler[*application.DeleteCharacterCommand, *application.DeleteCharacterHandlerResponse](deleteCharacterHandler))

	hE(mediatr.RegisterRequestHandler[*application.CreateLocationCommand, *application.CreateLocationHandlerResponse](createLocationHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateLocationCommand, *application.UpdateLocationHandlerResponse](updateLocationHandler))
	hE(mediatr.RegisterRequestHandler[*application.DeleteLocationCommand, *application.DeleteLocationHandlerResponse](deleteLocationHandler))

	hE(mediatr.RegisterRequestHandler[*application.UpdateMetadataCommand, *application.UpdateMetadataHandlerResponse](updateMetadataHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
