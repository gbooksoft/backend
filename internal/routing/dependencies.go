package routing

import (
	"github.com/gbooksoft/backend/internal/action"
	"github.com/gbooksoft/backend/internal/database"
	"github.com/gbooksoft/backend/internal/domain/account"
	"github.com/gbooksoft/backend/internal/domain/authorization"
)

var actionAuthorize *action.Authorize
var actionAccount *action.Account

func initializeDependencies() {
	actionAuthorize = action.NewAuthorize(
		authorization.NewService(database.GetDBPool()),
	)
	actionAccount = action.NewAccount(
		account.NewService(database.GetDBPool()),
	)
}
