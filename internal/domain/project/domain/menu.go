package domain

import "github.com/google/uuid"

type MenuItemType string

const (
	Chapter   MenuItemType = "chapter"
	Character MenuItemType = "character"
	Location  MenuItemType = "location"
)

type MenuId uuid.UUID

func (id MenuId) String() string {
	return uuid.UUID(id).String()
}

type MenuItem struct {
	Id        MenuId
	ProjectId ProjectId
	Title     string
	Type      MenuItemType
	Children  []MenuId
}

func NewDefaultMenu(projectId ProjectId) []*MenuItem {
	return []*MenuItem{
		{
			Id:        MenuId(uuid.New()),
			ProjectId: projectId,
			Title:     "Главы",
			Type:      Chapter,
		},
		{
			Id:        MenuId(uuid.New()),
			ProjectId: projectId,
			Title:     "Персонажи",
			Type:      Character,
		},
		{
			Id:        MenuId(uuid.New()),
			ProjectId: projectId,
			Title:     "Локации",
			Type:      Location,
		},
	}
}
