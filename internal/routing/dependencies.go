package routing

import (
	"github.com/gbooksoft/backend/internal/action"
	"github.com/gbooksoft/backend/internal/database"
	"github.com/gbooksoft/backend/internal/domain/account"
	"github.com/gbooksoft/backend/internal/domain/authorization"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/internal/domain/email/queue"
)

var actionAuthorize *action.Authorize
var actionAccount *action.Account

func initializeDependencies() {
	dbPool := database.GetDBPool()
	emailFacade := email.NewFacade(queue.NewEmailQueue(dbPool), queue.NewTemplateGenerator())

	actionAuthorize = action.NewAuthorize(
		authorization.NewService(dbPool, emailFacade),
	)
	actionAccount = action.NewAccount(
		account.NewService(dbPool, emailFacade),
	)
}
