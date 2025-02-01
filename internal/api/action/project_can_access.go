package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectCanAccess struct {
}

func NewProjectCanAccess() *ProjectCanAccess {
	return new(ProjectCanAccess)
}

func (a *ProjectCanAccess) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	projectId := ctx.Params.ByName("id")

	cmd := &application.CanAccessCommand{
		ProjectId: projectId,
		UserId:    claims.Subject,
	}

	resp, err := mediatr.Send[*application.CanAccessCommand, *application.CanAccessHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	if !resp.CanAccess {
		ctx.JSON(http.StatusForbidden, 1)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
