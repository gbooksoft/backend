package domain

import (
	"github.com/gbooksoft/backend/internal/gerr"
	"github.com/google/uuid"
)

type MenuItemType string

const (
	Chapter   MenuItemType = "chapter"
	Character MenuItemType = "character"
	Location  MenuItemType = "location"
)

type MenuItemId uuid.UUID

func (id MenuItemId) String() string {
	return uuid.UUID(id).String()
}

type Menu struct {
	Items []*MenuItem
}

func (m *Menu) findChildrenWithType(itemType MenuItemType) []*MenuItem {
	var result []*MenuItem
	for _, menuItem := range m.Items {
		if menuItem.Type == itemType && menuItem.ParentId != nil {
			result = append(result, menuItem)
		}
	}

	return result
}

func (m *Menu) AddNewItem(newMenuItemId MenuItemId, parentItemId MenuItemId) error {
	_, ok := m.findItemById(newMenuItemId)
	if ok {
		return gerr.NewBusinessLogicErr("Menu item already exists", gerr.ProjectMenuItemAlreadyExists, nil)
	}

	item, ok := m.findItemById(parentItemId)

	if !ok {
		return gerr.NewBusinessLogicErr("Could not find parent item", gerr.ProjectMenuParentItemNotFound, nil)
	}

	newItem := &MenuItem{
		Id:        newMenuItemId,
		ProjectId: item.ProjectId,
		Title:     "Unknown",
		Type:      item.Type,
		ParentId:  &parentItemId,
	}

	m.Items = append(m.Items, newItem)

	return nil
}

func (m *Menu) findItemById(id MenuItemId) (*MenuItem, bool) {
	for _, item := range m.Items {
		if item.Id == id {
			return item, true
		}
	}

	return nil, false
}

type MenuItem struct {
	Id        MenuItemId
	ProjectId ProjectId
	Title     string
	Type      MenuItemType
	ParentId  *MenuItemId
}

func NewDefaultMenu(projectId ProjectId) *Menu {
	chapterParentId := MenuItemId(uuid.New())
	return &Menu{
		Items: []*MenuItem{
			{
				Id:        chapterParentId,
				ProjectId: projectId,
				Title:     "Главы",
				Type:      Chapter,
			},
			{
				Id:        MenuItemId(uuid.New()),
				ProjectId: projectId,
				Title:     "Персонажи",
				Type:      Character,
			},
			{
				Id:        MenuItemId(uuid.New()),
				ProjectId: projectId,
				Title:     "Локации",
				Type:      Location,
			},
			{
				Id:        MenuItemId(uuid.New()),
				ProjectId: projectId,
				Title:     "Глава 1",
				Type:      Chapter,
				ParentId:  &chapterParentId,
			},
		},
	}
}
