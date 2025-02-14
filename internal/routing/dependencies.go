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
var projectGetLocationsAction *action.ProjectGetLocations
var projectUpdateLocationAction *action.ProjectUpdateLocation
var projectGetAction *action.ProjectGet
var projectUpdateAction *action.ProjectUpdate
var projectReorderMenuAction *action.ProjectReorderMenu
var projectCreateChapterAction *action.ProjectCreateChapter

var projectUpdateChapterAction *action.ProjectUpdateChapter
var projectDeleteChapterAction *action.ProjectDeleteChapter

var projectCreateCharacterAction *action.ProjectCreateCharacter
var projectUpdateCharacterAction *action.ProjectUpdateCharacter

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
	projectUpdateAction = action.NewProjectUpdate()
	projectReorderMenuAction = action.NewProjectReorderMenu()

	projectCreateChapterAction = action.NewProjectCreateChapter()
	projectUpdateChapterAction = action.NewProjectUpdateChapter()
	projectDeleteChapterAction = action.NewProjectDeleteChapter()

	projectCreateCharacterAction = action.NewProjectCreateCharacter()

	// genre
	genreGetAllAction = action.NewGenreGetAll()
}
