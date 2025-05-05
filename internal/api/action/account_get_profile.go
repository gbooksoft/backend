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

type accountProfileSectionResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	CoverURL  string `json:"cover_url"`
}
type accountProfileSecurityResponse struct {
	Email string `json:"email"`
}
type accountProfileSocialsResponse struct {
}

type accountProfileResponse struct {
	Id       string                         `json:"id"`
	Profile  accountProfileSectionResponse  `json:"profile"`
	Security accountProfileSecurityResponse `json:"security"`
	Socials  accountProfileSocialsResponse  `json:"socials,omitempty"`
}

func (a *AccountGetProfile) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(gerr.NewBusinessLogicErr("Unauthorized", gerr.AuthorizationUserNotFound))
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

	coverUrl, err := filesystem.GetPresignedURL(resp.CoverPath)

	if err != nil {
		// todo: default avatar?
		coverUrl = ""
	}

	ctx.JSON(http.StatusOK, accountProfileResponse{
		Id: claims.Subject,
		Profile: accountProfileSectionResponse{
			FirstName: resp.FirstName,
			LastName:  resp.LastName,
			Username:  resp.Username,
			AvatarURL: url,
			CoverURL:  coverUrl,
		},
		Security: accountProfileSecurityResponse{
			Email: resp.Email,
		},
		Socials: accountProfileSocialsResponse{},
	})
}
