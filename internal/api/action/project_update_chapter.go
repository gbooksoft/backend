package action

import (
	"github.com/gbooksoft/backend/internal/domain/project2/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gbooksoft/backend/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectUpdateChapter struct {
}

func NewProjectUpdateChapter() *ProjectUpdateChapter {
	return new(ProjectUpdateChapter)
}

type updateProjectChapterRequest struct {
	Title   string `json:"title" validate:"required,max=255"`
	Content string `json:"content" validate:"required"`
}

func (a *ProjectUpdateChapter) Invoke(ctx *gin.Context) {
	var req updateProjectChapterRequest
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
	chapterInternalId := ctx.Params.ByName("chapter_internal_id")

	cmd := &application.UpdateChapterCommand{
		ProjectId:         projectId,
		AuthorId:          claims.Subject,
		ChapterInternalId: chapterInternalId,
		Content:           req.Content,
		Title:             req.Title,
	}

	_, err = mediatr.Send[*application.UpdateChapterCommand, *application.UpdateChapterHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
