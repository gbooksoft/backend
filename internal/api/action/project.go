package action

import (
	"github.com/gbooksoft/backend/internal/api/action/response"
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Project struct {
	service *application.Service
}

func NewProject(service *application.Service) *Project {
	return &Project{
		service: service,
	}
}

type updateLocationRequest struct {
	Name        string  `json:"name" validate:"required,max=255,min=1"`
	Description *string `json:"description" validate:"max=5000,omitempty"`
}

func (a *Project) UpdateLocation(ctx *gin.Context) {
	var req updateLocationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(req)

	if err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	projectId := ctx.Params.ByName("id")
	menuItemId := ctx.Params.ByName("item_id")

	err = a.service.UpdateLocation(claims.Subject, projectId, menuItemId, req.Name, req.Description)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
	return
}

func (a *Project) GetProject(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	projectId := ctx.Params.ByName("id")

	project, err := a.service.GetProject(claims.Subject, projectId)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Struct(&projectItemResponse{
		Id:          project.Id,
		Title:       project.Title,
		Description: project.Description,
	}).Write()
}

type updateProjectRequest struct {
	Title       string  `json:"title" validate:"required,max=255,min=1"`
	Description *string `json:"description" validate:"max=255,omitempty"`
}

func (a *Project) UpdateProject(ctx *gin.Context) {
	var req updateProjectRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(req)

	if err != nil {
		response.Json(ctx.Writer).ErrValidation(err).Write()
		return
	}

	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	projectId := ctx.Params.ByName("id")

	err = a.service.UpdateProject(claims.Subject, projectId, req.Title, req.Description)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
	return
}
