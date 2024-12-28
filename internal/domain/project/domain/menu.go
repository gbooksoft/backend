package domain

import (
	"fmt"
	"github.com/gbooksoft/backend/internal/gerr"
	"github.com/google/uuid"
	"slices"
)

type MenuItemType string

const (
	TypeChapter            MenuItemType = "chapter"
	TypeCharacter          MenuItemType = "character"
	TypeLocation           MenuItemType = "location"
	TypeAdditionalElements MenuItemType = "additional_elements"
	TypeForeword           MenuItemType = "foreword"
	TypeAfterword          MenuItemType = "afterword"
	TypeMetadata           MenuItemType = "metadata"
)

type MenuItemId uuid.UUID

func (id MenuItemId) String() string {
	return uuid.UUID(id).String()
}

type Menu struct {
	Items []*MenuItem
}

func (m *Menu) findItemsWithType(itemTypes ...MenuItemType) []*MenuItem {
	var result []*MenuItem
	for _, menuItem := range m.Items {
		if slices.Contains(itemTypes, menuItem.Type) && menuItem.ParentId != nil {
			result = append(result, menuItem)
		}
	}

	return result
}

func (m *Menu) AddNewItem(newMenuItemId MenuItemId, parentItemId MenuItemId) (*MenuItem, error) {
	_, ok := m.findItemById(newMenuItemId)
	if ok {
		return nil, gerr.NewBusinessLogicErr("Menu item already exists", gerr.ProjectMenuItemAlreadyExists, nil)
	}

	item, ok := m.findItemById(parentItemId)

	if !ok {
		return nil, gerr.NewBusinessLogicErr("Could not find parent item", gerr.ProjectMenuParentItemNotFound, nil)
	}

	var title string

	if item.Type == TypeChapter {
		title = "Глава"
	} else if item.Type == TypeCharacter {
		title = "Персонаж"
	} else if item.Type == TypeLocation {
		title = "Локация"
	} else {
		return nil, gerr.NewBusinessLogicErr(fmt.Sprintf("can not add menu item to %s type", item.Type), gerr.ProjectCanNotAddMenuItemToType, nil)
	}

	newItem := &MenuItem{
		Id:        newMenuItemId,
		ProjectId: item.ProjectId,
		Title:     title,
		Type:      item.Type,
		ParentId:  &parentItemId,
		Order:     m.findLastOrderByParentId(parentItemId) + 1,
	}

	m.Items = append(m.Items, newItem)

	return newItem, nil
}

func (m *Menu) findLastOrderByParentId(id MenuItemId) int {
	maxOrder := 0
	for _, item := range m.Items {
		if item.ParentId != nil && *item.ParentId == id {
			maxOrder = max(maxOrder, item.Order)
		}
	}

	return maxOrder
}

func (m *Menu) findItemById(id MenuItemId) (*MenuItem, bool) {
	for _, item := range m.Items {
		if item.Id == id {
			return item, true
		}
	}

	return nil, false
}

func (m *Menu) UpdateItemTitle(itemUuid MenuItemId, title string) error {
	item, found := m.findItemById(itemUuid)

	if !found {
		return gerr.NewBusinessLogicErr("Could not find parent item", gerr.ProjectMenuItemNotFound, nil)
	}

	item.Title = title

	return nil
}

func (m *Menu) SelectItem(itemUuid MenuItemId) {
	for _, item := range m.Items {
		if item.Id == itemUuid {
			item.Selected = true
		} else {
			item.Selected = false
		}
	}
}

type MenuItem struct {
	Id        MenuItemId
	ProjectId ProjectId
	Title     string
	Type      MenuItemType
	ParentId  *MenuItemId
	Selected  bool
	Order     int
}

func NewDefaultMenu(projectId ProjectId) (*Menu, MenuItemId) {
	chapterParentId := MenuItemId(uuid.New())
	additionalElementsParentId := MenuItemId(uuid.New())

	return &Menu{
		Items: []*MenuItem{
			{
				Id:        chapterParentId,
				ProjectId: projectId,
				Title:     "Главы",
				Type:      TypeChapter,
				Order:     1,
			},
			{
				Id:        MenuItemId(uuid.New()),
				ProjectId: projectId,
				Title:     "Персонажи",
				Type:      TypeCharacter,
				Order:     2,
			},
			{
				Id:        MenuItemId(uuid.New()),
				ProjectId: projectId,
				Title:     "Локации",
				Type:      TypeLocation,
				Order:     3,
			},
			{
				Id:        MenuItemId(uuid.New()),
				ProjectId: projectId,
				Title:     "Глава 1",
				Type:      TypeChapter,
				ParentId:  &chapterParentId,
				Selected:  true,
				Order:     1,
			},
			{
				Id:        additionalElementsParentId,
				ProjectId: projectId,
				Title:     "Элементы книги",
				Type:      TypeAdditionalElements,
				Order:     4,
			},
			{
				Id:        MenuItemId(uuid.New()),
				ProjectId: projectId,
				ParentId:  &additionalElementsParentId,
				Title:     "Предисловие",
				Type:      TypeForeword,
				Order:     1,
			},
			{
				Id:        MenuItemId(uuid.New()),
				ProjectId: projectId,
				ParentId:  &additionalElementsParentId,
				Title:     "Послесловие",
				Type:      TypeAfterword,
				Order:     2,
			},
			{
				Id:        MenuItemId(uuid.New()),
				ProjectId: projectId,
				ParentId:  &additionalElementsParentId,
				Title:     "Метаданные",
				Type:      TypeMetadata,
				Order:     3,
			},
		},
	}, chapterParentId
}
