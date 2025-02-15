package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectUpdate struct {
}

func NewProjectUpdate() *ProjectUpdate {
	return new(ProjectUpdate)
}

type updateProjectRequest struct {
	Title       string   `json:"title" validate:"required,max=255,min=1"`
	Description *string  `json:"description" validate:"omitempty,max=255"`
	Genres      []string `json:"genres" validate:"required,min=1,max=5,dive,uuid4"`
}

func (a *ProjectUpdate) Invoke(ctx *gin.Context) {
	var req updateProjectRequest
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

	cmd := &application.UpdateProjectCommand{
		ProjectId:   projectId,
		Title:       req.Title,
		Description: req.Description,
		AuthorId:    claims.Subject,
		Genres:      req.Genres,
	}

	_, err = mediatr.Send[*application.UpdateProjectCommand, *application.UpdateProjectHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
