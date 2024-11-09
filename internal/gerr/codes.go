package gerr

type ErrorCode int64

func (c *ErrorCode) ToInt() int64 {
	return int64(*c)
}

const (
	_ ErrorCode = 1000 + iota
	AuthorizationUserNotFound
	AuthorizationUserAlreadyExists
	AuthorizationRefreshTokenExpired
	AuthorizationRefreshTokenNotFound
)

const (
	_ ErrorCode = 2000 + iota
)
