package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
)

type ProjectGetEditors struct {
}

func NewProjectGetEditors() *ProjectGetEditors {
	return new(ProjectGetEditors)
}

type getEditorItemResponse struct {
	MenuItemId  string `json:"menu_item_id"`
	EditorState string `json:"editor_state"`
}

func (a *ProjectGetEditors) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	projectId := ctx.Params.ByName("id")

	query := &application.GetEditorsQuery{
		ProjectId: projectId,
		UserId:    claims.Subject,
	}

	resp, err := mediatr.Send[*application.GetEditorsQuery, *application.GetEditorsHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	var result []getEditorItemResponse
	for _, item := range resp.Editors {
		result = append(result, getEditorItemResponse{
			MenuItemId:  item.MenuItemId,
			EditorState: item.EditorState,
		})
	}

	ctx.JSON(200, gin.H{
		"items": result,
	})
	return
}
