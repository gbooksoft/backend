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

type projectGenreItemResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type projectItemResponse struct {
	Id          string                     `json:"id"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	Genres      []projectGenreItemResponse `json:"genres"`
}

func (a *ProjectGetAll) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	query := &application.GetAllProjectsQuery{
		AuthorId: claims.Subject,
	}

	resp, err := mediatr.Send[*application.GetAllProjectsQuery, *application.GetAllProjectsHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	var result []projectItemResponse
	for _, item := range resp.Projects {
		responseItem := projectItemResponse{
			Id:          item.Id,
			Title:       item.Title,
			Description: item.Description,
		}

		for _, genre := range item.Genres {
			responseItem.Genres = append(responseItem.Genres, projectGenreItemResponse{
				Id:   genre.Id,
				Name: genre.Name,
			})
		}

		result = append(result, responseItem)
	}

	ctx.JSON(http.StatusOK, getAllResponse{Items: result})
	return
}
