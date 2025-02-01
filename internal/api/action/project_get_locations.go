package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectGetLocations struct {
}

func NewProjectGetLocations() *ProjectGetLocations {
	return new(ProjectGetLocations)
}

type getLocationItemResponse struct {
	MenuItemId  string `json:"menu_item_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (a *ProjectGetLocations) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	projectId := ctx.Params.ByName("id")

	query := &application.GetLocationsQuery{
		ProjectId: projectId,
		UserId:    claims.Subject,
	}
	resp, err := mediatr.Send[*application.GetLocationsQuery, *application.GetLocationsHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	var result []getLocationItemResponse
	for _, item := range resp.Locations {
		result = append(result, getLocationItemResponse{
			MenuItemId:  item.MenuItemId,
			Name:        item.Name,
			Description: item.Description,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items": result,
	})
	return
}
