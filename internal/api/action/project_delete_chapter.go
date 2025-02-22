package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectDeleteChapter struct {
}

func NewProjectDeleteChapter() *ProjectDeleteChapter {
	return new(ProjectDeleteChapter)
}

func (a *ProjectDeleteChapter) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	projectId := ctx.Params.ByName("id")
	chapterInternalId := ctx.Params.ByName("chapter_internal_id")

	cmd := &application.DeleteChapterCommand{
		ProjectId:         projectId,
		AuthorId:          claims.Subject,
		ChapterInternalId: chapterInternalId,
	}

	_, err = mediatr.Send[*application.DeleteChapterCommand, *application.DeleteChapterHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
