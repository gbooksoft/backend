package action

import (
	"github.com/gbooksoft/backend/internal/domain/project2/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectDeleteLocation struct {
}

func NewProjectDeleteLocation() *ProjectDeleteLocation {
	return new(ProjectDeleteLocation)
}

func (a *ProjectDeleteLocation) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	projectId := ctx.Params.ByName("id")
	locationInternalId := ctx.Params.ByName("location_internal_id")

	cmd := &application.DeleteLocationCommand{
		ProjectId:          projectId,
		AuthorId:           claims.Subject,
		LocationInternalId: locationInternalId,
	}

	_, err = mediatr.Send[*application.DeleteLocationCommand, *application.DeleteLocationHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
