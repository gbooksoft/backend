package routing

import (
	"github.com/gbooksoft/backend/internal/action"
	"github.com/gbooksoft/backend/internal/database"
	"github.com/gbooksoft/backend/internal/domain/account"
	"github.com/gbooksoft/backend/internal/domain/authorization"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/internal/domain/email/queue"
	project_application "github.com/gbooksoft/backend/internal/domain/project/application"
	project_infrastructure "github.com/gbooksoft/backend/internal/domain/project/infrastructure"
)

var actionAuthorize *action.Authorize
var actionAccount *action.Account
var actionProject *action.Project

func initializeDependencies() {
	dbPool := database.GetDBPool()
	emailFacade := email.NewFacade(queue.NewEmailQueue(dbPool), queue.NewTemplateGenerator())

	actionAuthorize = action.NewAuthorize(
		authorization.NewService(dbPool, emailFacade),
	)
	actionAccount = action.NewAccount(
		account.NewService(dbPool, emailFacade),
	)

	actionProject = action.NewProject(
		project_application.NewService(
			project_infrastructure.NewMySqlProjectRepository(dbPool),
			project_infrastructure.NewMySqlMenuRepository(dbPool),
		),
	)
}
