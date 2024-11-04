package routing

import (
	"github.com/gbooksoft/backend/internal/action"
	"github.com/gbooksoft/backend/internal/database"
	"github.com/gbooksoft/backend/internal/domain/authorization"
)

var actionAuthorize *action.Authorize

func initializeDependencies() {
	actionAuthorize = action.NewAuthorize(
		authorization.NewService(database.GetDBPool()),
	)
}
