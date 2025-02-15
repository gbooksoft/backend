package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gbooksoft/backend/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectCreateCharacter struct {
}

func NewProjectCreateCharacter() *ProjectCreateCharacter {
	return new(ProjectCreateCharacter)
}

type createProjectCharacterRequest struct {
	Id   string `json:"id" validate:"required,uuid"`
	Name string `json:"name" validate:"required,max=255"`
}

func (a *ProjectCreateCharacter) Invoke(ctx *gin.Context) {
	var req createProjectCharacterRequest
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

	cmd := &application.CreateCharacterCommand{
		Id:        req.Id,
		ProjectId: projectId,
		AuthorId:  claims.Subject,
		Name:      req.Name,
	}

	_, err = mediatr.Send[*application.CreateCharacterCommand, *application.CreateCharacterHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
