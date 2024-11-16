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

	err = s.menuRepository.AddMenu(tx, menu)

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
