package action

import (
	"github.com/gbooksoft/backend/internal/domain/project2/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectUpdateCharacter struct {
}

func NewProjectUpdateCharacter() *ProjectUpdateCharacter {
	return new(ProjectUpdateCharacter)
}

type updateCharacterRequest struct {
	Name        string                            `json:"name" validate:"required,max=255,min=1"`
	Description *string                           `json:"description" validate:"max=5000,omitempty"`
	Attributes  []updateCharacterAttributeRequest `json:"attributes"`
}
type updateCharacterAttributeRequest struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func (a *ProjectUpdateCharacter) Invoke(ctx *gin.Context) {
	var req updateCharacterRequest
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
	characterInternalId := ctx.Params.ByName("character_internal_id")

	var charAttrs map[string]string
	if len(req.Attributes) > 0 {
		charAttrs = make(map[string]string, len(req.Attributes))
		for _, attr := range req.Attributes {
			charAttrs[attr.Field] = attr.Value
		}
	}

	cmd := &application.UpdateCharacterCommand{
		CharacterInternalId: characterInternalId,
		AuthorId:            claims.Subject,
		ProjectId:           projectId,
		Name:                req.Name,
		Description:         req.Description,
		Attributes:          charAttrs,
	}

	_, err = mediatr.Send[*application.UpdateCharacterCommand, *application.UpdateCharacterHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
	return
}
