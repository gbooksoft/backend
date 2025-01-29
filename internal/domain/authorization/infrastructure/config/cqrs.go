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

	hE(mediatr.RegisterRequestHandler[*application.SignInCommand, *application.SignInHandlerResponse](signInHandler))
	hE(mediatr.RegisterRequestHandler[*application.SignUpCommand, *application.SignUpHandlerResponse](signUpHandler))
	hE(mediatr.RegisterRequestHandler[*application.RefreshCommand, *application.RefreshHandlerResponse](refreshHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
