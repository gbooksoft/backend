package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectSaveEditorState struct {
}

func NewProjectSaveEditorState() *ProjectSaveEditorState {
	return new(ProjectSaveEditorState)
}

type saveEditorStateRequest struct {
	EditorState string `json:"editor_state" validate:"required"`
}

func (a *ProjectSaveEditorState) Invoke(ctx *gin.Context) {
	var req saveEditorStateRequest
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
	menuItemId := ctx.Params.ByName("item_id")

	cmd := &application.SaveEditorStateCommand{
		UserId:      claims.Subject,
		ProjectId:   projectId,
		MenuItemId:  menuItemId,
		EditorState: req.EditorState,
	}
	_, err = mediatr.Send[*application.SaveEditorStateCommand, *application.SaveEditorStateHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, 1)
}
