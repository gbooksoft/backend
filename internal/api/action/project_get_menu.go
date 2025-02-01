package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectGetMenu struct {
}

func NewProjectGetMenu() *ProjectGetMenu {
	return &ProjectGetMenu{}
}

func (a *ProjectGetMenu) Invoke(ctx *gin.Context) {
	projectId := ctx.Params.ByName("id")
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	query := &application.GetMenuQuery{
		UserId:    claims.Subject,
		ProjectId: projectId,
	}
	resp, err := mediatr.Send[*application.GetMenuQuery, *application.GetMenuHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	responseItems := a.groupMenuByParentId(resp.Menu)

	expandItemsQuery := &application.GetExpandedItemsQuery{
		UserId:    claims.Subject,
		ProjectId: projectId,
	}

	expandItemsResp, err := mediatr.Send[*application.GetExpandedItemsQuery, *application.GetExpandedItemsHandlerResponse](ctx, expandItemsQuery)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"menu":           responseItems,
		"expanded_items": expandItemsResp.ItemsIds,
	})
	return
}

type getProjectMenuItemResponse struct {
	Id         string                       `json:"id"`
	Title      string                       `json:"title"`
	Type       string                       `json:"type"`
	ParentId   *string                      `json:"parent_id"`
	Children   []getProjectMenuItemResponse `json:"children"`
	IsSelected bool                         `json:"is_selected"`
	Order      int                          `json:"order"`
}

func (a *ProjectGetMenu) groupMenuByParentId(menu *readmodel.Menu) []getProjectMenuItemResponse {
	grouped := make(map[string][]getProjectMenuItemResponse)
	for _, item := range menu.Items {
		var parentId string

		if item.ParentId == nil {
			parentId = "ROOT"
		} else {
			parentId = *item.ParentId
		}

		menuItem := getProjectMenuItemResponse{
			Id:         item.Id,
			Title:      item.Title,
			Type:       item.ItemType,
			ParentId:   item.ParentId,
			Children:   nil,
			IsSelected: item.IsSelected,
			Order:      item.Order,
		}

		grouped[parentId] = append(grouped[parentId], menuItem)
	}

	for parentId, children := range grouped {
		for childI := range children {
			curId := grouped[parentId][childI].Id

			if itsChildren, ok := grouped[curId]; ok {
				grouped[parentId][childI].Children = itsChildren
			}
		}
	}

	if v, ok := grouped["ROOT"]; ok {
		return v
	}

	return []getProjectMenuItemResponse{}
}
