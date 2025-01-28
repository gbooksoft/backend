package action

import (
	"github.com/gbooksoft/backend/internal/domain/account/application"
	"github.com/gbooksoft/backend/pkg/filesystem"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type AccountGetProfile struct {
}

func NewAccountGetProfile() *AccountGetProfile {
	return new(AccountGetProfile)
}

func (a *AccountGetProfile) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(gerr.NewBusinessLogicErr("Unauthorized", gerr.AuthorizationUserNotFound, nil))
		return
	}

	query := &application.GetProfileQuery{
		UserId: claims.Subject,
	}

	resp, err := mediatr.Send[*application.GetProfileQuery, *application.GetProfileHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	url, err := filesystem.GetPresignedURL(resp.AvatarPath)

	if err != nil {
		// todo: default avatar?
		url = ""
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":    claims.Subject,
		"email":      resp.Email,
		"first_name": resp.FirstName,
		"last_name":  resp.LastName,
		"avatar_url": url,
	})
}
