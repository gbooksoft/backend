package routing

import (
	"github.com/gbooksoft/backend/internal/api/action"
	project_application "github.com/gbooksoft/backend/internal/domain/project/application"
	project_infrastructure "github.com/gbooksoft/backend/internal/domain/project/infrastructure"
	"github.com/gbooksoft/backend/pkg/database"
)

var actionProject *action.Project

var accountGetProfileAction *action.AccountGetProfile
var accountUpdatePasswordAction *action.AccountUpdatePassword
var accountUpdateSettingsAction *action.AccountUpdateSettings
var accountConfirmEmailChangeAction *action.AccountConfirmEmailChange

var authSignInAction *action.AuthSignIn
var authSignUpAction *action.AuthSignUp
var authRefreshAction *action.AuthRefresh
var authLogoutAction *action.AuthLogout
var authConfirmRegistrationAction *action.AuthConfirmRegistration
var authForgotPasswordAction *action.AuthForgotPassword
var authResetPasswordAction *action.AuthResetPassword

var projectCreateAction *action.ProjectCreate
var projectGetAllAction *action.ProjectGetAll
var projectGetRecentAction *action.ProjectGetRecent
var projectCanAccessAction *action.ProjectCanAccess
var projectGetMetadataAction *action.ProjectGetMetadata
var projectUpdateMetadataAction *action.ProjectUpdateMetadata
var projectGetMenuAction *action.ProjectGetMenu
var projectAddMenuItemAction *action.ProjectAddMenuItem
var projectToggleExpandMenuItemAction *action.ProjectToggleExpandMenuItem
var projectUpdateMenuItemAction *action.ProjectUpdateMenuItem
var projectDeleteMenuItemAction *action.ProjectDeleteMenuItem
var projectGetEditorsAction *action.ProjectGetEditors
var projectSaveEditorStateAction *action.ProjectSaveEditorState
var projectGetCharactersAction *action.ProjectGetCharacters
var projectUpdateCharacterAction *action.ProjectUpdateCharacter
var projectGetLocationsAction *action.ProjectGetLocations
var projectUpdateLocationAction *action.ProjectUpdateLocation
var projectGetAction *action.ProjectGet

func initializeDependencies() {
	dbPool := database.GetDBPool()

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
	authLogoutAction = action.NewAuthLogout()
	authConfirmRegistrationAction = action.NewAuthConfirmRegistration()
	authForgotPasswordAction = action.NewAuthForgotPassword()
	authResetPasswordAction = action.NewAuthResetPassword()

	// project
	projectCreateAction = action.NewProjectCreate()
	projectGetAllAction = action.NewProjectGetAll()
	projectGetRecentAction = action.NewProjectGetRecent()
	projectCanAccessAction = action.NewProjectCanAccess()
	projectGetMetadataAction = action.NewProjectGetMetadata()
	projectUpdateMetadataAction = action.NewProjectUpdateMetadata()
	projectGetMenuAction = action.NewProjectGetMenu()
	projectAddMenuItemAction = action.NewProjectAddMenuItem()
	projectToggleExpandMenuItemAction = action.NewProjectToggleExpandMenuItem()
	projectUpdateMenuItemAction = action.NewProjectUpdateMenuItem()
	projectDeleteMenuItemAction = action.NewProjectDeleteMenuItem()
	projectGetEditorsAction = action.NewProjectGetEditors()
	projectSaveEditorStateAction = action.NewProjectSaveEditorState()
	projectGetCharactersAction = action.NewProjectGetCharacters()
	projectUpdateCharacterAction = action.NewProjectUpdateCharacter()
	projectGetLocationsAction = action.NewProjectGetLocations()
	projectUpdateLocationAction = action.NewProjectUpdateLocation()
	projectGetAction = action.NewProjectGet()
}
