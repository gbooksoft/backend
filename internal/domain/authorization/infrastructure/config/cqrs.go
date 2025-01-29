package config

import (
	"github.com/gbooksoft/backend/internal/domain/authorization/application"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/internal/domain/email/queue"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/mehdihadeli/go-mediatr"
)

func RegisterMediatr() {
	dbPool := database.GetDBPool()
	emailFacade := email.NewFacade(queue.NewEmailQueue(dbPool), queue.NewTemplateGenerator())

	signInHandler := application.NewSignInHandler(dbPool)
	signUpHandler := application.NewSignUpHandler(dbPool, emailFacade)
	refreshHandler := application.NewRefreshHandler(dbPool)
	logoutHandler := application.NewLogoutHandler(dbPool)
	confirmRegistrationHandler := application.NewConfirmRegistrationHandler(dbPool)

	hE(mediatr.RegisterRequestHandler[*application.SignInCommand, *application.SignInHandlerResponse](signInHandler))
	hE(mediatr.RegisterRequestHandler[*application.SignUpCommand, *application.SignUpHandlerResponse](signUpHandler))
	hE(mediatr.RegisterRequestHandler[*application.RefreshCommand, *application.RefreshHandlerResponse](refreshHandler))
	hE(mediatr.RegisterRequestHandler[*application.LogoutCommand, *application.LogoutHandlerResponse](logoutHandler))
	hE(mediatr.RegisterRequestHandler[*application.ConfirmRegistrationCommand, *application.ConfirmRegistrationHandlerResponse](confirmRegistrationHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
