package action

import (
	"github.com/gbooksoft/backend/internal/domain/project2/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectReorderChapter struct {
}

func NewProjectReorderChapter() *ProjectReorderChapter {
	return new(ProjectReorderChapter)
}

type reorderRequest struct {
	Order int `json:"order" validate:"required"`
}

func (a *ProjectReorderChapter) Invoke(ctx *gin.Context) {
	var req reorderRequest
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
	chapterInternalId := ctx.Params.ByName("chapter_internal_id")

	cmd := &application.ReorderChapterCommand{
		AuthorId:          claims.Subject,
		ProjectId:         projectId,
		Order:             req.Order,
		ChapterInternalId: chapterInternalId,
	}

	_, err = mediatr.Send[*application.ReorderChapterCommand, *application.ReorderChapterHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
