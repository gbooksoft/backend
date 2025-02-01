package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectGetRecent struct {
}

func NewProjectGetRecent() *ProjectGetRecent {
	return new(ProjectGetRecent)
}
func (a *ProjectGetRecent) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	query := &application.GetRecentProjectsQuery{
		UserId: claims.Subject,
		Limit:  3,
	}

	resp, err := mediatr.Send[*application.GetRecentProjectsQuery, *application.GetRecentProjectsHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	var result []projectItemResponse
	for _, item := range resp.Projects {
		result = append(result, projectItemResponse{
			Id:          item.Id,
			Title:       item.Title,
			Description: item.Description,
		})
	}

	ctx.JSON(http.StatusOK, getAllResponse{Items: result})
	return
}
