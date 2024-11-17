package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
	"github.com/gbooksoft/backend/internal/gerr"
	"github.com/google/uuid"
)

type Service struct {
	projectRepository domain.ProjectRepository
	menuRepository    domain.MenuRepository
	menuReadModel     readmodel.MenuReadModel
}

func NewService(
	projectRepository domain.ProjectRepository,
	menuRepository domain.MenuRepository,
	menuReadModel readmodel.MenuReadModel,
) *Service {
	return &Service{
		projectRepository: projectRepository,
		menuRepository:    menuRepository,
		menuReadModel:     menuReadModel,
	}
}

func (s *Service) CreateProject(userId string, title string) (string, error) {
	tx, err := s.projectRepository.BeginTransaction(context.TODO())

	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	userUuid := domain.UserId(uuid.MustParse(userId))

	exists, err := s.projectRepository.ExistsByUserIdAndTitle(tx, userUuid, title)

	if err != nil {
		return "", err
	}

	if exists {
		return "", gerr.NewBusinessLogicErr("Project with this title already exists", gerr.ProjectAlreadyExists, nil)
	}

	projectId := domain.ProjectId(s.projectRepository.GetNextId())

	project := domain.NewProject(projectId, userUuid, title)

	err = s.projectRepository.Add(tx, project)

	if err != nil {
		return "", err
	}

	menu := domain.NewDefaultMenu(projectId)

	err = s.menuRepository.AddMenu(tx, menu, nil)

	if err != nil {
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return projectId.String(), nil
}

func (s *Service) GetProjectMenu(userId string, projectId string) (*readmodel.Menu, error) {
	menu, err := s.menuReadModel.GetByUserIdAndProjectId(userId, projectId)

	if err != nil {
		return nil, err
	}

	if menu == nil || len(menu.Items) == 0 {
		return nil, gerr.NewBusinessLogicErr("Menu not found", gerr.NotFound, nil)
	}

	return menu, nil
}

func (s *Service) AddMenuItem(newMenuItemId, projectId, parentId string) error {
	projectUUid := domain.ProjectId(uuid.MustParse(projectId))
	newMenuItemUUid := domain.MenuItemId(uuid.MustParse(newMenuItemId))
	parentUUid := domain.MenuItemId(uuid.MustParse(parentId))

	tx, err := s.menuRepository.BeginTransaction(context.TODO())

	if err != nil {
		return err
	}
	defer tx.Rollback()

	menu, err := s.menuRepository.FindByProjectId(tx, projectUUid)

	if err != nil {
		return err
	}

	err = menu.AddNewItem(newMenuItemUUid, parentUUid)

	if err != nil {
		return err
	}

	err = s.menuRepository.Save(tx, projectUUid, menu)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetExpandedItems(userId string, projectId string) ([]string, error) {
	return s.menuReadModel.GetExpandedItems(userId, projectId)
}

func (s *Service) ExpandItem(userId string, projectId string, menuItemId string) error {
	err := s.menuReadModel.ExpandItem(userId, projectId, menuItemId)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CollapseItem(userId string, projectId string, menuItemId string) error {
	err := s.menuReadModel.CollapseItem(userId, projectId, menuItemId)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UserCanAccessToProject(userId string, projectId string) (bool, error) {
	return s.projectRepository.UserHasAccess(domain.UserId(uuid.MustParse(userId)), domain.ProjectId(uuid.MustParse(projectId)))
}
