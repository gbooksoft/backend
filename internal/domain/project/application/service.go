package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/internal/gerr"
	"github.com/google/uuid"
)

type Service struct {
	projectRepository domain.ProjectRepository
	menuRepository    domain.MenuRepository
}

func NewService(projectRepository domain.ProjectRepository, menuRepository domain.MenuRepository) *Service {
	return &Service{
		projectRepository: projectRepository,
		menuRepository:    menuRepository,
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
