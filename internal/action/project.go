package action

import (
	"github.com/gbooksoft/backend/internal/action/response"
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
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

type createProjectRequest struct {
	Title string `json:"title" binding:"required" validate:"required,max=255"`
}

type createProjectResponse struct {
	Id string `json:"id"`
}

func (a *Project) Create(ctx *gin.Context) {
	var req createProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	projectId, err := a.service.CreateProject(claims.Subject, req.Title)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Struct(createProjectResponse{
		Id: projectId,
	}).Write()
	return
}

type getProjectMenuResponse struct {
	Menu []getProjectMenuItemResponse `json:"menu"`
}

type getProjectMenuItemResponse struct {
	Id       string                       `json:"id"`
	Title    string                       `json:"title"`
	Type     string                       `json:"type"`
	ParentId *string                      `json:"parent_id"`
	Children []getProjectMenuItemResponse `json:"children"`
}

func (a *Project) GetProjectMenu(ctx *gin.Context) {
	projectId := ctx.Params.ByName("id")
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	menu, err := a.service.GetProjectMenu(claims.Subject, projectId)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	responseItems := a.groupMenuByParentId(menu)

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Struct(getProjectMenuResponse{
		Menu: responseItems,
	}).Write()

	return
}

func (a *Project) groupMenuByParentId(menu *readmodel.Menu) []getProjectMenuItemResponse {
	grouped := make(map[string][]getProjectMenuItemResponse)
	for _, item := range menu.Items {
		var parentId string

		if item.ParentId == nil {
			parentId = "ROOT"
		} else {
			parentId = *item.ParentId
		}

		grouped[parentId] = append(grouped[parentId], getProjectMenuItemResponse{
			Id:       item.Id,
			Title:    item.Title,
			Type:     item.ItemType,
			ParentId: item.ParentId,
			Children: nil})
	}

	for parentId, children := range grouped {
		for childI, _ := range children {
			curId := grouped[parentId][childI].Id

			if itsChildren, ok := grouped[curId]; ok {
				grouped[parentId][childI].Children = itsChildren
			}
		}
	}

	if v, ok := grouped["ROOT"]; ok {
		return v
	}

	return []getProjectMenuItemResponse{}
}
