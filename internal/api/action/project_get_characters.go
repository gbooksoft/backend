package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectGetCharacters struct {
}

func NewProjectGetCharacters() *ProjectGetCharacters {
	return new(ProjectGetCharacters)
}

type getCharacterItemResponse struct {
	MenuItemId  string                              `json:"menu_item_id"`
	Name        string                              `json:"name"`
	Description string                              `json:"description"`
	Attributes  []getCharacterItemAttributeResponse `json:"attributes"`
}
type getCharacterItemAttributeResponse struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func (a *ProjectGetCharacters) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	projectId := ctx.Params.ByName("id")

	query := &application.GetCharactersQuery{
		UserId:    claims.Subject,
		ProjectId: projectId,
	}

	resp, err := mediatr.Send[*application.GetCharactersQuery, *application.GetCharactersHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	var result []getCharacterItemResponse
	for _, item := range resp.Characters {
		var attributes []getCharacterItemAttributeResponse

		for _, attr := range item.Attributes {
			attributes = append(attributes, getCharacterItemAttributeResponse{
				Field: attr.Field,
				Value: attr.Value,
			})
		}

		result = append(result, getCharacterItemResponse{
			MenuItemId:  item.MenuItemId,
			Name:        item.Name,
			Description: item.Description,
			Attributes:  attributes,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items": result,
	})
	return
}
