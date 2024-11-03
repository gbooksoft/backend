package authorization

import (
	"github.com/gbooksoft/backend/internal/database"
	"github.com/go-oauth2/mysql/v4"
)

func newTokenStore() *mysql.Store {
	return mysql.NewStoreWithDB(database.GetDBPool(), "oauth2_token", 0)
}
