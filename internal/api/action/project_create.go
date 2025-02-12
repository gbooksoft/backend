package action

import (
	"github.com/gbooksoft/backend/internal/domain/project2/application"
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
	Genres      []string `json:"genres" validate:"min=1,max=5,dive,uuid4"`
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
		AuthorId:     claims.Subject,
		UserFullName: userFullName,
		Title:        req.Title,
		Description:  req.Description,
		Genres:       req.Genres,
	}

	resp, err := mediatr.Send[*application.CreateProjectCommand, *application.CreateProjectHandlerResponse](ctx, cmd)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"project_id": resp.ProjectId,
	})
	return
}
