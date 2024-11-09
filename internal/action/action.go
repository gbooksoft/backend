package action

import (
	"errors"
	"github.com/gbooksoft/backend/internal/domain/authorization"
	"github.com/gin-gonic/gin"
)

func getJwtClaims(ctx *gin.Context) (*authorization.CustomClaims, error) {
	token, exists := ctx.Get("jwt_token_claims")
	if !exists {
		return nil, errors.New("token does not exist")
	}

	jwtClaims, ok := token.(*authorization.CustomClaims)
	if !ok {
		return nil, errors.New("token is not correct type")
	}

	return jwtClaims, nil
}
