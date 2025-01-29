package config

import (
	"github.com/gbooksoft/backend/internal/domain/authorization/application"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/mehdihadeli/go-mediatr"
)

func RegisterMediatr() {
	dbPool := database.GetDBPool()

	signInHandler := application.NewSignInHandler(dbPool)

	hE(mediatr.RegisterRequestHandler[*application.SignInCommand, *application.SignInHandlerResponse](signInHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
