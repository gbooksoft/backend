package routing

import (
	"github.com/gbooksoft/backend/internal/api/action"
)

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
var projectGetAction *action.ProjectGet
var projectUpdateAction *action.ProjectUpdate

var projectCreateChapterAction *action.ProjectCreateChapter
var projectUpdateChapterAction *action.ProjectUpdateChapter
var projectDeleteChapterAction *action.ProjectDeleteChapter
var projectReorderChapterAction *action.ProjectReorderChapter

var projectCreateCharacterAction *action.ProjectCreateCharacter
var projectUpdateCharacterAction *action.ProjectUpdateCharacter
var projectDeleteCharacterAction *action.ProjectDeleteCharacter

var projectCreateLocationAction *action.ProjectCreateLocation
var projectUpdateLocationAction *action.ProjectUpdateLocation
var projectDeleteLocationAction *action.ProjectDeleteLocation

var projectUpdateMetadataAction *action.ProjectUpdateMetadata

var projectPublishAction *action.ProjectPublish
var projectRepublishAction *action.ProjectRepublish

var genreGetAllAction *action.GenreGetAll

func initializeDependencies() {
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
	projectGetAction = action.NewProjectGet()
	projectUpdateAction = action.NewProjectUpdate()

	projectCreateChapterAction = action.NewProjectCreateChapter()
	projectUpdateChapterAction = action.NewProjectUpdateChapter()
	projectDeleteChapterAction = action.NewProjectDeleteChapter()
	projectReorderChapterAction = action.NewProjectReorderChapter()

	projectCreateCharacterAction = action.NewProjectCreateCharacter()
	projectUpdateCharacterAction = action.NewProjectUpdateCharacter()
	projectDeleteCharacterAction = action.NewProjectDeleteCharacter()

	projectCreateLocationAction = action.NewProjectCreateLocation()
	projectUpdateLocationAction = action.NewProjectUpdateLocation()
	projectDeleteLocationAction = action.NewProjectDeleteLocation()

	projectUpdateMetadataAction = action.NewProjectUpdateMetadata()

	projectPublishAction = action.NewProjectPublish()
	projectRepublishAction = action.NewProjectRepublish()

	// genre
	genreGetAllAction = action.NewGenreGetAll()
}
