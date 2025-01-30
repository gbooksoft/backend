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
	sendForgotPasswordNotificationHandler := application.NewSendForgotPasswordNotificationHandler(dbPool, emailFacade)
	resetPasswordHandler := application.NewResetPasswordHandler(dbPool, emailFacade)

	hE(mediatr.RegisterRequestHandler[*application.SignInCommand, *application.SignInHandlerResponse](signInHandler))
	hE(mediatr.RegisterRequestHandler[*application.SignUpCommand, *application.SignUpHandlerResponse](signUpHandler))
	hE(mediatr.RegisterRequestHandler[*application.RefreshCommand, *application.RefreshHandlerResponse](refreshHandler))
	hE(mediatr.RegisterRequestHandler[*application.LogoutCommand, *application.LogoutHandlerResponse](logoutHandler))
	hE(mediatr.RegisterRequestHandler[*application.ConfirmRegistrationCommand, *application.ConfirmRegistrationHandlerResponse](confirmRegistrationHandler))
	hE(mediatr.RegisterRequestHandler[*application.SendForgotPasswordNotificationCommand, *application.SendForgotPasswordNotificationHandlerResponse](sendForgotPasswordNotificationHandler))
	hE(mediatr.RegisterRequestHandler[*application.ResetPasswordCommand, *application.ResetPasswordHandlerResponse](resetPasswordHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
