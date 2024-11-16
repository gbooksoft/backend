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
		Children:  nil,
	}

	item.Children = append(item.Children, newItem)

	return nil
}

func (m *Menu) findItemById(id MenuItemId) (*MenuItem, bool) {
	return m.findIn(m.Items, id)
}

func (m *Menu) findIn(items []*MenuItem, id MenuItemId) (*MenuItem, bool) {
	for _, item := range items {
		if item.Id == id {
			return item, true
		}

		if len(item.Children) > 0 {
			itemInChildren, found := m.findIn(item.Children, id)

			if found {
				return itemInChildren, true
			}
		}
	}
	return nil, false
}

type MenuItem struct {
	Id        MenuItemId
	ProjectId ProjectId
	Title     string
	Type      MenuItemType
	Children  []*MenuItem
}

func NewDefaultMenu(projectId ProjectId) *Menu {
	return &Menu{
		Items: []*MenuItem{
			{
				Id:        MenuItemId(uuid.New()),
				ProjectId: projectId,
				Title:     "Главы",
				Type:      Chapter,
				Children: []*MenuItem{
					{
						Id:        MenuItemId(uuid.New()),
						ProjectId: projectId,
						Title:     "Глава 1",
						Type:      Chapter,
					},
				},
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
		},
	}
}
