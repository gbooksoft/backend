package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectDeleteCharacter struct {
}

func NewProjectDeleteCharacter() *ProjectDeleteCharacter {
	return new(ProjectDeleteCharacter)
}

func (a *ProjectDeleteCharacter) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	projectId := ctx.Params.ByName("id")
	characterInternalId := ctx.Params.ByName("character_internal_id")

	cmd := &application.DeleteCharacterCommand{
		ProjectId:           projectId,
		AuthorId:            claims.Subject,
		CharacterInternalId: characterInternalId,
	}

	_, err = mediatr.Send[*application.DeleteCharacterCommand, *application.DeleteCharacterHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
