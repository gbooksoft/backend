package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectGet struct {
}

func NewProjectGet() *ProjectGet {
	return new(ProjectGet)
}

func (a *ProjectGet) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	projectId := ctx.Params.ByName("id")

	query := &application.GetProjectQuery{
		ProjectId: projectId,
		UserId:    claims.Subject,
	}

	resp, err := mediatr.Send[*application.GetProjectQuery, *application.GetProjectHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          resp.Project.Id,
		"title":       resp.Project.Title,
		"description": resp.Project.Description,
	})
}
