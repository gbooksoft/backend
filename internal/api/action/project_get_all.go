package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectGetAll struct {
}

func NewProjectGetAll() *ProjectGetAll {
	return new(ProjectGetAll)
}

type getAllResponse struct {
	Items []projectItemResponse `json:"projects"`
}
type projectItemResponse struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (a *ProjectGetAll) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	query := &application.GetAllProjectsQuery{
		UserId: claims.Subject,
	}

	resp, err := mediatr.Send[*application.GetAllProjectsQuery, *application.GetAllProjectsQueryHandlerResponse](ctx, query)

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
