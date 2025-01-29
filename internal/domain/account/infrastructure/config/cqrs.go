package config

import (
	"github.com/gbooksoft/backend/internal/domain/account/application"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/internal/domain/email/queue"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/mehdihadeli/go-mediatr"
)

func RegisterMediatr() {
	dbPool := database.GetDBPool()
	emailFacade := email.NewFacade(queue.NewEmailQueue(dbPool), queue.NewTemplateGenerator())

	getProfileHandler := application.NewGetProfileHandler(dbPool)
	updatePasswordHandler := application.NewUpdatePasswordHandler(dbPool, emailFacade)
	updateSettingsHandler := application.NewUpdateSettingsHandler(dbPool, emailFacade)
	uploadAvatarHandler := application.NewUploadAvatarHandler(dbPool)
	confirmEmailChangeHandler := application.NewConfirmEmailChangeCommandHandler(dbPool)

	hE(mediatr.RegisterRequestHandler[*application.GetProfileQuery, *application.GetProfileHandlerResponse](getProfileHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdatePasswordCommand, *application.UpdatePasswordHandlerResponse](updatePasswordHandler))
	hE(mediatr.RegisterRequestHandler[*application.UpdateSettingsCommand, *application.UpdateSettingsHandlerResponse](updateSettingsHandler))
	hE(mediatr.RegisterRequestHandler[*application.UploadAvatarCommand, *application.UploadAvatarHandlerResponse](uploadAvatarHandler))
	hE(mediatr.RegisterRequestHandler[*application.ConfirmEmailChangeCommand, *application.ConfirmEmailChangeCommandHandlerResponse](confirmEmailChangeHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
