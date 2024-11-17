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
	Menu          []getProjectMenuItemResponse `json:"menu"`
	ExpandedItems []string                     `json:"expanded_items"`
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
	expandedItems, _ := a.service.GetExpandedItems(claims.Subject, projectId)

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Struct(getProjectMenuResponse{
		Menu:          responseItems,
		ExpandedItems: expandedItems,
	}).Write()

	return
}

type addMenuItemRequest struct {
	Id       string `json:"id" validate:"required,uuid"`
	ParentId string `json:"parent_id" validate:"required,uuid"`
}

func (a *Project) AddMenuItem(ctx *gin.Context) {
	var req addMenuItemRequest
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

	projectId := ctx.Params.ByName("id")

	err = a.service.AddMenuItem(req.Id, projectId, req.ParentId)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
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

type toggleExpandRequest struct {
	Expand bool `json:"expand" validate:"required"`
}

func (a *Project) ToggleExpandMenuItem(ctx *gin.Context) {
	var req toggleExpandRequest
	if err := ctx.ShouldBind(&req); err != nil {
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

	if req.Expand {
		err = a.service.ExpandItem(claims.Subject, projectId, menuItemId)
	} else {
		err = a.service.CollapseItem(claims.Subject, projectId, menuItemId)
	}

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
	return
}

func (a *Project) CanAccess(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}
	projectId := ctx.Params.ByName("id")

	has, err := a.service.UserCanAccessToProject(claims.Subject, projectId)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	if !has {
		response.Json(ctx.Writer).StatusCode(http.StatusForbidden).Write()
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
	return
}
