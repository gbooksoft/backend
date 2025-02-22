package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectUpdateLocation struct {
}

func NewProjectUpdateLocation() *ProjectUpdateLocation {
	return new(ProjectUpdateLocation)
}

type updateLocationRequest struct {
	Name        string  `json:"name" validate:"required,max=255,min=1"`
	Description *string `json:"description" validate:"omitempty,max=5000"`
}

func (a *ProjectUpdateLocation) Invoke(ctx *gin.Context) {
	var req updateLocationRequest
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
	locationInternalId := ctx.Params.ByName("location_internal_id")

	cmd := &application.UpdateLocationCommand{
		AuthorId:           claims.Subject,
		ProjectId:          projectId,
		LocationInternalId: locationInternalId,
		Name:               req.Name,
		Description:        req.Description,
	}

	_, err = mediatr.Send[*application.UpdateLocationCommand, *application.UpdateLocationHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
