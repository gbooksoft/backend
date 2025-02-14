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

	hE(mediatr.RegisterRequestHandler[*application.CreateProjectCommand, *application.CreateProjectHandlerResponse](createProjectHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateProjectCommand, *application.UpdateProjectHandlerResponse](updateProjectHandler))
	hE(mediatr.RegisterRequestHandler[*application.CreateChapterCommand, *application.CreateChapterHandlerResponse](createChapterHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateChapterCommand, *application.UpdateChapterHandlerResponse](updateChapterHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
