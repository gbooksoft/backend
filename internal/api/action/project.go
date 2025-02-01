package action

import (
	"github.com/gbooksoft/backend/internal/api/action/response"
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

type getProjectMenuResponse struct {
	Menu          []getProjectMenuItemResponse `json:"menu"`
	ExpandedItems []string                     `json:"expanded_items"`
}

type getProjectMenuItemResponse struct {
	Id         string                       `json:"id"`
	Title      string                       `json:"title"`
	Type       string                       `json:"type"`
	ParentId   *string                      `json:"parent_id"`
	Children   []getProjectMenuItemResponse `json:"children"`
	IsSelected bool                         `json:"is_selected"`
	Order      int                          `json:"order"`
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

	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	projectId := ctx.Params.ByName("id")

	err = a.service.AddMenuItem(req.Id, claims.Subject, projectId, req.ParentId)

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

		menuItem := getProjectMenuItemResponse{
			Id:         item.Id,
			Title:      item.Title,
			Type:       item.ItemType,
			ParentId:   item.ParentId,
			Children:   nil,
			IsSelected: item.IsSelected,
			Order:      item.Order,
		}

		grouped[parentId] = append(grouped[parentId], menuItem)
	}

	for parentId, children := range grouped {
		for childI := range children {
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

type metadataRequest struct {
	Author   string `json:"author,omitempty" validate:"omitempty,max=255"`
	Genre    string `json:"genre,omitempty" validate:"omitempty,max=255"`
	Keywords string `json:"keywords,omitempty" validate:"omitempty,max=255"`
	Language string `json:"language,omitempty" validate:"omitempty,max=255"`
	Isbn     string `json:"isbn,omitempty" validate:"omitempty,max=255"`
}

func (a *Project) CreateOrUpdateMetadata(ctx *gin.Context) {
	var req metadataRequest
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

	err = a.service.CreateOrUpdateMetadata(claims.Subject, projectId, application.MetadataFields{
		Author:   req.Author,
		Genre:    req.Genre,
		Keywords: req.Keywords,
		Language: req.Language,
		Isbn:     req.Isbn,
	})

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
}

type updateMenuItemRequest struct {
	Title    string `json:"title" validate:"max=255"`
	Selected bool   `json:"selected" validate:"boolean"`
}

func (a *Project) UpdateMenuItem(ctx *gin.Context) {
	var req updateMenuItemRequest
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

	err = a.service.UpdateMenuItem(claims.Subject, projectId, menuItemId, req.Title, req.Selected)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
}

func (a *Project) DeleteMenuItem(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	projectId := ctx.Params.ByName("id")
	menuItemId := ctx.Params.ByName("item_id")

	err = a.service.DeleteMenuItem(claims.Subject, projectId, menuItemId)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
}

type saveEditorStateRequest struct {
	EditorState string `json:"editor_state" validate:"required"`
}

func (a *Project) SaveEditorState(ctx *gin.Context) {
	var req saveEditorStateRequest
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

	err = a.service.SaveEditorState(claims.Subject, projectId, menuItemId, req.EditorState)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
}

type getEditorsResponse struct {
	Items []getEditorItemResponse `json:"items"`
}
type getEditorItemResponse struct {
	MenuItemId  string `json:"menu_item_id"`
	EditorState string `json:"editor_state"`
}

func (a *Project) GetEditors(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	projectId := ctx.Params.ByName("id")

	editors, err := a.service.FetchEditors(claims.Subject, projectId)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	var result []getEditorItemResponse
	for _, item := range editors {
		result = append(result, getEditorItemResponse{
			MenuItemId:  item.MenuItemId,
			EditorState: item.EditorState,
		})
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Struct(getEditorsResponse{
		Items: result,
	}).Write()
	return
}

type getCharactersResponse struct {
	Items []getCharacterItemResponse `json:"items"`
}
type getCharacterItemResponse struct {
	MenuItemId  string                              `json:"menu_item_id"`
	Name        string                              `json:"name"`
	Description string                              `json:"description"`
	Attributes  []getCharacterItemAttributeResponse `json:"attributes"`
}
type getCharacterItemAttributeResponse struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func (a *Project) GetCharacters(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	projectId := ctx.Params.ByName("id")

	characters, err := a.service.FetchCharacters(claims.Subject, projectId)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	var result []getCharacterItemResponse
	for _, item := range characters {
		var attributes []getCharacterItemAttributeResponse

		for _, attr := range item.Attributes {
			attributes = append(attributes, getCharacterItemAttributeResponse{
				Field: attr.Field,
				Value: attr.Value,
			})
		}

		result = append(result, getCharacterItemResponse{
			MenuItemId:  item.MenuItemId,
			Name:        item.Name,
			Description: item.Description,
			Attributes:  attributes,
		})
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Struct(getCharactersResponse{
		Items: result,
	}).Write()
	return
}

type getLocationsResponse struct {
	Items []getLocationItemResponse `json:"items"`
}
type getLocationItemResponse struct {
	MenuItemId  string `json:"menu_item_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (a *Project) GetLocations(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		response.Json(ctx.Writer).Unauthorized(err.Error()).Write()
		return
	}

	projectId := ctx.Params.ByName("id")

	characters, err := a.service.FetchLocations(claims.Subject, projectId)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	var result []getLocationItemResponse
	for _, item := range characters {

		result = append(result, getLocationItemResponse{
			MenuItemId:  item.MenuItemId,
			Name:        item.Name,
			Description: item.Description,
		})
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Struct(getLocationsResponse{
		Items: result,
	}).Write()
	return
}

type updateCharacterAttributeRequest struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type updateCharacterRequest struct {
	Name        string                            `json:"name" validate:"required,max=255,min=1"`
	Description *string                           `json:"description" validate:"max=5000,omitempty"`
	Attributes  []updateCharacterAttributeRequest `json:"attributes"`
}

func (a *Project) UpdateCharacter(ctx *gin.Context) {
	var req updateCharacterRequest
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

	var charAttrs []application.CharacterAttribute
	for _, attr := range req.Attributes {
		charAttrs = append(charAttrs, application.CharacterAttribute{
			Field: attr.Field,
			Value: attr.Value,
		})
	}
	err = a.service.UpdateCharacter(claims.Subject, projectId, menuItemId, req.Name, req.Description, charAttrs)

	if err != nil {
		response.HandleErr(err, ctx.Writer)
		return
	}

	response.Json(ctx.Writer).StatusCode(http.StatusOK).Write()
	return
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
