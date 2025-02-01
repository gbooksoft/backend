package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectDeleteMenuItem struct {
}

func NewProjectDeleteMenuItem() *ProjectDeleteMenuItem {
	return new(ProjectDeleteMenuItem)
}

func (a *ProjectDeleteMenuItem) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	projectId := ctx.Params.ByName("id")
	menuItemId := ctx.Params.ByName("item_id")

	cmd := &application.DeleteMenuItemCommand{
		UserId:     claims.Subject,
		ProjectId:  projectId,
		MenuItemId: menuItemId,
	}

	_, err = mediatr.Send[*application.DeleteMenuItemCommand, *application.DeleteMenuItemHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
}
