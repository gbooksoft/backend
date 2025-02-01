package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectUpdateMenuItem struct {
}

func NewProjectUpdateMenuItem() *ProjectUpdateMenuItem {
	return new(ProjectUpdateMenuItem)
}

type updateMenuItemRequest struct {
	Title    string `json:"title" validate:"max=255"`
	Selected bool   `json:"selected" validate:"boolean"`
}

func (a *ProjectUpdateMenuItem) Invoke(ctx *gin.Context) {
	var req updateMenuItemRequest
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
	menuItemId := ctx.Params.ByName("item_id")

	cmd := &application.UpdateMenuItemCommand{
		UserId:      claims.Subject,
		ProjectId:   projectId,
		MenuItemId:  menuItemId,
		Title:       req.Title,
		SetSelected: req.Selected,
	}

	_, err = mediatr.Send[*application.UpdateMenuItemCommand, *application.UpdateMenuItemHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
}
