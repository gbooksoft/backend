package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectGetPublication struct {
}

func NewProjectGetPublication() *ProjectGetPublication {
	return new(ProjectGetPublication)
}

type publicationChapter struct {
	Id      string `json:"project_chapter_id"`
	Title   string `json:"title"`
	Version int    `json:"version"`
	Order   int    `json:"order"`
}
type publicationResponse struct {
	IsPublished bool                 `json:"is_published"`
	Chapters    []publicationChapter `json:"chapters"`
}

func (a *ProjectGetPublication) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	projectId := ctx.Params.ByName("id")

	query := &application.GetPublicationQuery{
		ProjectId: projectId,
		AuthorId:  claims.Subject,
	}

	resp, err := mediatr.Send[*application.GetPublicationQuery, *application.GetPublicationHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	if resp == nil {
		ctx.JSON(http.StatusOK, publicationResponse{
			IsPublished: false,
			Chapters:    []publicationChapter{},
		})
		return
	}

	var chapters []publicationChapter
	for _, c := range resp.Publication.Chapters {
		chapters = append(chapters, publicationChapter{
			Id:      c.ProjectChapterId,
			Title:   c.Title,
			Version: c.Version,
			Order:   c.Order,
		})
	}

	ctx.JSON(http.StatusOK, publicationResponse{
		IsPublished: true,
		Chapters:    chapters,
	})
}
