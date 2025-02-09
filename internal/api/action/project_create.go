package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gbooksoft/backend/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectCreate struct {
}

func NewProjectCreate() *ProjectCreate {
	return new(ProjectCreate)
}

type createProjectRequest struct {
	Title       string   `json:"title" binding:"required" validate:"required,max=255"`
	Description string   `json:"description" validate:"max=255"`
	Genres      []string `json:"genres" validate:"uuid4_slice,max_slice_size=5,min_slice_size=1"`
	Tags        []string `json:"tags" validate:"max_slice_size=5"`
}

func (a *ProjectCreate) Invoke(ctx *gin.Context) {
	var req createProjectRequest
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

	userFullName := claims.FirstName + " " + claims.LastName

	cmd := &application.CreateProjectCommand{
		UserId:       claims.Subject,
		UserFullName: userFullName,
		Title:        req.Title,
		Description:  req.Description,
		Genres:       req.Genres,
		Tags:         req.Tags,
	}

	resp, err := mediatr.Send[*application.CreateProjectCommand, *application.CreateProjectHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	expandCmd := &application.ExpandItemCommand{
		MenuItemId: resp.ItemIdToExpand,
		ProjectId:  resp.ProjectId,
		UserId:     claims.Subject,
	}

	_, err = mediatr.Send[*application.ExpandItemCommand, *application.ExpandItemHandlerResponse](ctx, expandCmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"project_id": resp.ProjectId,
	})
	return
}
