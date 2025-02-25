package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gbooksoft/backend/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectRepublish struct {
}

func NewProjectRepublish() *ProjectRepublish {
	return new(ProjectRepublish)
}

type republishProjectRequest struct {
	Chapters []string `json:"chapters" validate:"omitempty"`
}

func (a *ProjectRepublish) Invoke(ctx *gin.Context) {
	var req republishProjectRequest
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

	cmd := &application.RepublishProjectCommand{
		ProjectId: projectId,
		AuthorId:  claims.Subject,
		Chapters:  req.Chapters,
	}

	_, err = mediatr.Send[*application.RepublishProjectCommand, *application.RepublishProjectHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
