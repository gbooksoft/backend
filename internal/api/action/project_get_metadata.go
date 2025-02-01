package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectGetMetadata struct {
}

func NewProjectGetMetadata() *ProjectGetMetadata {
	return new(ProjectGetMetadata)
}
func (a *ProjectGetMetadata) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	projectId := ctx.Params.ByName("id")

	query := &application.GetMetadataQuery{
		UserId:    claims.Subject,
		ProjectId: projectId,
	}

	resp, err := mediatr.Send[*application.GetMetadataQuery, *application.GetMetadataHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"author":   resp.Metadata.Author,
		"genre":    resp.Metadata.Genre,
		"keywords": resp.Metadata.Keywords,
		"language": resp.Metadata.Language,
		"isbn":     resp.Metadata.Isbn,
	})
}
