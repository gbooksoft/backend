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
			project_infrastructure.NewMySqlMenuReadModel(dbPool),
			project_infrastructure.NewMySqlEditorRepository(dbPool),
			project_infrastructure.NewMySqlEditorReadModel(dbPool),
			project_infrastructure.NewMySqlProjectReadModel(dbPool),
			project_infrastructure.NewMySqlCharacterRepository(dbPool),
			project_infrastructure.NewMySqlCharacterReadModel(dbPool),
			project_infrastructure.NewMySqlLocationRepository(dbPool),
			project_infrastructure.NewMySqlLocationReadModel(dbPool),
		),
	)
}
