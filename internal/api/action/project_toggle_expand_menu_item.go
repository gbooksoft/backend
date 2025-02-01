package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectToggleExpandMenuItem struct {
}

func NewProjectToggleExpandMenuItem() *ProjectToggleExpandMenuItem {
	return new(ProjectToggleExpandMenuItem)
}

type toggleExpandRequest struct {
	Expand bool `json:"expand" validate:"required"`
}

func (a *ProjectToggleExpandMenuItem) Invoke(ctx *gin.Context) {
	var req toggleExpandRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	projectId := ctx.Params.ByName("id")
	menuItemId := ctx.Params.ByName("item_id")

	if req.Expand {
		cmd := &application.ExpandItemCommand{
			UserId:     claims.Subject,
			ProjectId:  projectId,
			MenuItemId: menuItemId,
		}

		_, err = mediatr.Send[*application.ExpandItemCommand, *application.ExpandItemHandlerResponse](ctx, cmd)
	} else {
		cmd := &application.CollapseItemCommand{
			UserId:     claims.Subject,
			ProjectId:  projectId,
			MenuItemId: menuItemId,
		}

		_, err = mediatr.Send[*application.CollapseItemCommand, *application.CollapseItemHandlerResponse](ctx, cmd)
	}

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
