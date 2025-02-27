package config

import (
	"github.com/gbooksoft/backend/internal/domain/author/application"
	"github.com/gbooksoft/backend/internal/domain/author/infrastructure"
	"github.com/gbooksoft/backend/pkg/database"
	"github.com/mehdihadeli/go-mediatr"
)

func RegisterMediatr() {
	authorReadModel := infrastructure.NewMySqlAuthorReadModel(database.GetGormDB())
	getAuthorByUsername := application.NewGetAuthorByUsernameHandler(authorReadModel)

	hE(mediatr.RegisterRequestHandler[*application.GetAuthorByUsernameQuery, *application.GetAuthorByUsernameHandlerResponse](getAuthorByUsername))
}

func hE(err error) {
	if err != nil {
		panic(err)
	}
}
