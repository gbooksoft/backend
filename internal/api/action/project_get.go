package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectGet struct {
}

func NewProjectGet() *ProjectGet {
	return new(ProjectGet)
}

func (a *ProjectGet) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	projectId := ctx.Params.ByName("id")

	query := &application.GetProjectQuery{
		ProjectId: projectId,
		UserId:    claims.Subject,
	}

	resp, err := mediatr.Send[*application.GetProjectQuery, *application.GetProjectHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	responseItem := projectItemResponse{
		Id:          resp.Project.Id,
		Title:       resp.Project.Title,
		Description: resp.Project.Description,
	}

	for _, genre := range resp.Project.Genres {
		responseItem.Genres = append(responseItem.Genres, projectGenreItemResponse{
			Id:   genre.Id,
			Name: genre.Name,
		})
	}

	ctx.JSON(http.StatusOK, responseItem)
}
