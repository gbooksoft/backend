package gerr

type ErrorCode int64

func (c *ErrorCode) ToInt() int64 {
	return int64(*c)
}

const (
	_ ErrorCode = 0 + iota
	NotFound
)

const (
	_ ErrorCode = 1000 + iota
	AuthorizationUserNotFound
	AuthorizationUserAlreadyExists
	AuthorizationRefreshTokenExpired
	AuthorizationRefreshTokenNotFound
	AuthorizationConfirmationNoFound
	AuthorizationResetPasswordRequestNotFound
)

const (
	_ ErrorCode = 2000 + iota
	AccountUserNotFound
	AccountNothingToUpdate
	AccountUserAlreadyExists
	AccountUserWithEmailAlreadyExists
	AccountEmailChangeRequestWithEmailAlreadyExists
	AccountEmailChangeRequestNotFound
	AccountEmailChangeRequestExpired
	AccountPasswordIsIncorrect
)

const (
	_ ErrorCode = 3000 + iota
	ProjectAlreadyExists
	ProjectMenuItemAlreadyExists
	ProjectMenuParentItemNotFound
	ProjectMenuItemNotFound
	ProjectCanNotAddMenuItemToType
	ProjectMetadataNotFound
	ProjectIncorrectAmountOfItems
	ProjectIncorrectOrderingOfItems
	ProjectShouldHaveMaxFiveGenres
	ProjectChapterAlreadyExists
	ProjectChapterNotFound
	ProjectCharacterNotFound
)

const (
	_ ErrorCode = 4000 + iota
	GenreNotExists
)
