package action

import (
	"github.com/gbooksoft/backend/internal/domain/author/application"
	"github.com/gbooksoft/backend/pkg/filesystem"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
)

type AuthorGetByUsername struct {
	Username string
}

func NewAuthorGetByUsername() *AuthorGetByUsername {
	return new(AuthorGetByUsername)
}

type authorResponse struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarUrl string `json:"avatar_url"`
}

func (a *AuthorGetByUsername) Invoke(ctx *gin.Context) {
	username := ctx.Params.ByName("username")

	query := &application.GetAuthorByUsernameQuery{
		Username: username,
	}

	resp, err := mediatr.Send[*application.GetAuthorByUsernameQuery, *application.GetAuthorByUsernameHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	url, err := filesystem.GetPresignedURL(resp.Author.AvatarPath)

	if err != nil {
		// todo: default avatar?
		url = ""
	}

	ctx.JSON(200, authorResponse{
		Id:        resp.Author.Id,
		Username:  resp.Author.Username,
		FirstName: resp.Author.FirstName,
		LastName:  resp.Author.LastName,
		AvatarUrl: url,
	})
}
