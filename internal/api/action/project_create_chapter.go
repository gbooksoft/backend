package action

import (
	"github.com/gbooksoft/backend/internal/domain/project2/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gbooksoft/backend/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectCreateChapter struct {
}

func NewProjectCreateChapter() *ProjectCreateChapter {
	return new(ProjectCreateChapter)
}

type createProjectChapterRequest struct {
	Id   string `json:"id" validate:"required,uuid"`
	Name string `json:"name" validate:"required,max=255"`
}

func (a *ProjectCreateChapter) Invoke(ctx *gin.Context) {
	var req createProjectChapterRequest
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

	cmd := &application.CreateChapterCommand{
		Id:        req.Id,
		ProjectId: projectId,
		AuthorId:  claims.Subject,
		Name:      req.Name,
	}

	_, err = mediatr.Send[*application.CreateChapterCommand, *application.CreateChapterHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
