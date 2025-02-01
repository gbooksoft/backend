package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectAddMenuItem struct {
}

func NewProjectAddMenuItem() *ProjectAddMenuItem {
	return new(ProjectAddMenuItem)
}

type addMenuItemRequest struct {
	Id       string `json:"id" validate:"required,uuid"`
	ParentId string `json:"parent_id" validate:"required,uuid"`
}

func (a *ProjectAddMenuItem) Invoke(ctx *gin.Context) {
	var req addMenuItemRequest
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

	cmd := &application.AddMenuItemCommand{
		NewItemId: req.Id,
		UserId:    claims.Subject,
		ProjectId: projectId,
		ParentId:  req.ParentId,
	}

	_, err = mediatr.Send[*application.AddMenuItemCommand, *application.AddMenuItemHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
}
