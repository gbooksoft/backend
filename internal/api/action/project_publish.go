package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gbooksoft/backend/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectPublish struct {
}

func NewProjectPublish() *ProjectPublish {
	return new(ProjectPublish)
}

type publishProjectRequest struct {
	Chapters []string `json:"chapters" validate:"omitempty"`
}

func (a *ProjectPublish) Invoke(ctx *gin.Context) {
	var req publishProjectRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(gerr.NewValidationError(err))
		return
	}

	validate := validator.NewCustomValidator()

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

	cmd := &application.PublishProjectCommand{
		ProjectId: projectId,
		AuthorId:  claims.Subject,
		Chapters:  req.Chapters,
	}

	_, err = mediatr.Send[*application.PublishProjectCommand, *application.PublishProjectHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
