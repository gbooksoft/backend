package config

import (
	"github.com/gbooksoft/backend/internal/domain/account/application"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/mehdihadeli/go-mediatr"
)

func RegisterMediatr() {
	getProfileHandler := application.NewGetProfileHandler(database.GetDBPool())

	hE(mediatr.RegisterRequestHandler[*application.GetProfileQuery, *application.GetProfileHandlerResponse](getProfileHandler))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
