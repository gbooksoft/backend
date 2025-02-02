package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectReorderMenu struct {
}

func NewProjectReorderMenu() *ProjectReorderMenu {
	return new(ProjectReorderMenu)
}

type reorderMenuRequest struct {
	MenuItems []menuItemRequest `json:"menu_items" validate:"required"`
}

type menuItemRequest struct {
	ItemId   string `json:"item_id" validate:"required"`
	ParentId string `json:"parent_id" validate:"required"`
	Order    int    `json:"order" validate:"required"`
}

func (a *ProjectReorderMenu) Invoke(ctx *gin.Context) {
	var req reorderMenuRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(req)

	if err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	projectId := ctx.Params.ByName("id")

	var items []application.MenuItem
	for _, menuItem := range req.MenuItems {
		items = append(items, application.MenuItem{
			ItemId:   menuItem.ItemId,
			ParentId: menuItem.ParentId,
			Order:    menuItem.Order,
		})
	}

	cmd := &application.ReorderMenuCommand{
		UserId:    claims.Subject,
		ProjectId: projectId,
		MenuItems: items,
	}

	_, err = mediatr.Send[*application.ReorderMenuCommand, *application.ReorderMenuHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
