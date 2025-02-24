package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectPublish struct {
}

func NewProjectPublish() *ProjectPublish {
	return new(ProjectPublish)
}

func (a *ProjectPublish) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	projectId := ctx.Params.ByName("id")

	cmd := &application.PublishProjectCommand{
		ProjectId: projectId,
		AuthorId:  claims.Subject,
	}

	_, err = mediatr.Send[*application.PublishProjectCommand, *application.PublishProjectHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
