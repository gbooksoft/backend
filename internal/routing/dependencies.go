package routing

import (
	"github.com/gbooksoft/backend/internal/api/action"
	"github.com/gbooksoft/backend/internal/domain/authorization"
	"github.com/gbooksoft/backend/internal/domain/email"
	"github.com/gbooksoft/backend/internal/domain/email/queue"
	project_application "github.com/gbooksoft/backend/internal/domain/project/application"
	project_infrastructure "github.com/gbooksoft/backend/internal/domain/project/infrastructure"
	"github.com/gbooksoft/backend/pkg/database"
)

var actionAuthorize *action.Authorize
var actionProject *action.Project

var accountGetProfileAction *action.AccountGetProfile
var accountUpdatePasswordAction *action.AccountUpdatePassword
var accountUpdateSettingsAction *action.AccountUpdateSettings
var accountConfirmEmailChangeAction *action.AccountConfirmEmailChange

var authSignInAction *action.AuthSignIn
var authSignUpAction *action.AuthSignUp
var authRefreshAction *action.AuthRefresh

func initializeDependencies() {
	dbPool := database.GetDBPool()
	emailFacade := email.NewFacade(queue.NewEmailQueue(dbPool), queue.NewTemplateGenerator())

	actionAuthorize = action.NewAuthorize(
		authorization.NewService(dbPool, emailFacade),
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

	// account
	accountGetProfileAction = action.NewAccountGetProfile()
	accountUpdatePasswordAction = action.NewAccountUpdatePassword()
	accountUpdateSettingsAction = action.NewAccountUpdateSettings()
	accountConfirmEmailChangeAction = action.NewAccountConfirmEmailChange()

	// auth
	authSignInAction = action.NewAuthSignIn()
	authSignUpAction = action.NewAuthSignUp()
	authRefreshAction = action.NewAuthRefresh()
}
