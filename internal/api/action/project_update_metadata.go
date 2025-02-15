package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectUpdateMetadata struct {
}

func NewProjectUpdateMetadata() *ProjectUpdateMetadata {
	return new(ProjectUpdateMetadata)
}

type metadataRequest struct {
	Author   string `json:"author,omitempty" validate:"omitempty,max=255"`
	Keywords string `json:"keywords,omitempty" validate:"omitempty,max=255"`
	Language string `json:"language,omitempty" validate:"omitempty,max=255"`
	Isbn     string `json:"isbn,omitempty" validate:"omitempty,max=255"`
}

func (a *ProjectUpdateMetadata) Invoke(ctx *gin.Context) {
	var req metadataRequest
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

	cmd := &application.UpdateMetadataCommand{
		ProjectId: projectId,
		AuthorId:  claims.Subject,
		Author:    req.Author,
		Keywords:  req.Keywords,
		Language:  req.Language,
		Isbn:      req.Isbn,
	}

	_, err = mediatr.Send[*application.UpdateMetadataCommand, *application.UpdateMetadataHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
}
