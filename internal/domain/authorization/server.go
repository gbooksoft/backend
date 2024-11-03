package authorization

import (
	"context"
	"errors"
	"github.com/go-oauth2/oauth2/v4"
	errors2 "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"time"
)

var invalidCredentials = errors.New("invalid_credential")

func CreateOAuthServer() *server.Server {
	manager := manage.NewDefaultManager()

	manager.MapTokenStorage(newTokenStore())
	manager.MapClientStorage(newClientStore())
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte(os.Getenv("JWT_KEY")), jwt.SigningMethodHS512))

	manager.SetPasswordTokenCfg(&manage.Config{
		AccessTokenExp:    time.Hour * 2,
		RefreshTokenExp:   time.Hour * 24 * 14,
		IsGenerateRefresh: true,
	})

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(false)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetAllowedGrantType(oauth2.PasswordCredentials)

	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		if username == "admin" && password == "admin" {
			return "admin", nil
		}

		return "", invalidCredentials
	})

	srv.SetInternalErrorHandler(func(err error) (re *errors2.Response) {
		if errors.Is(err, invalidCredentials) {
			return &errors2.Response{
				StatusCode:  http.StatusUnauthorized,
				Description: "Invalid username or password",
				Error:       invalidCredentials,
			}
		}

		return &errors2.Response{}
	})

	return srv
}
