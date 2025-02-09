package action

import (
	"github.com/gbooksoft/backend/internal/domain/genre/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type GenreGetAll struct {
}

func NewGenreGetAll() *GenreGetAll {
	return new(GenreGetAll)
}

type getAllGenresResponse struct {
	Items []genreItemResponse `json:"genres"`
}
type genreItemResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (a *GenreGetAll) Invoke(ctx *gin.Context) {
	query := &application.GetAllGenresQuery{
		Name: ctx.Query("name"),
	}

	resp, err := mediatr.Send[*application.GetAllGenresQuery, *application.GetAllGenresHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	var result []genreItemResponse
	for _, item := range resp.Genres {
		result = append(result, genreItemResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	ctx.JSON(http.StatusOK, getAllGenresResponse{Items: result})
	return
}
