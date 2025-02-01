package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type Service struct {
	projectRepository   domain.ProjectRepository
	menuRepository      domain.MenuRepository
	menuReadModel       readmodel.MenuReadModel
	editorRepository    domain.EditorRepository
	editorReadModel     readmodel.EditorReadModel
	projectReadModel    readmodel.ProjectReadModel
	characterRepository domain.CharacterRepository
	characterReadModel  readmodel.CharacterReadModel
	locationRepository  domain.LocationRepository
	locationReadModel   readmodel.LocationReadModel
}

func NewService(
	projectRepository domain.ProjectRepository,
	menuRepository domain.MenuRepository,
	menuReadModel readmodel.MenuReadModel,
	editorRepository domain.EditorRepository,
	editorReadModel readmodel.EditorReadModel,
	projectReadModel readmodel.ProjectReadModel,
	characterRepository domain.CharacterRepository,
	characterReadModel readmodel.CharacterReadModel,
	locationRepository domain.LocationRepository,
	locationReadModel readmodel.LocationReadModel,
) *Service {
	return &Service{
		projectRepository:   projectRepository,
		menuRepository:      menuRepository,
		menuReadModel:       menuReadModel,
		editorRepository:    editorRepository,
		editorReadModel:     editorReadModel,
		projectReadModel:    projectReadModel,
		characterRepository: characterRepository,
		characterReadModel:  characterReadModel,
		locationRepository:  locationRepository,
		locationReadModel:   locationReadModel,
	}
}

func (s *Service) FetchCharacters(userId string, projectId string) ([]readmodel.Character, error) {
	return s.characterReadModel.GetCharacters(userId, projectId)
}

func (s *Service) FetchLocations(userId string, projectId string) ([]readmodel.Location, error) {
	return s.locationReadModel.GetLocations(userId, projectId)
}

type CharacterAttribute struct {
	Field string
	Value string
}

func (s *Service) UpdateCharacter(
	userId string,
	projectId string,
	menuItemId string,
	name string,
	description *string,
	attributes []CharacterAttribute,
) error {
	userUUid := domain.UserId(uuid.MustParse(userId))
	projectUUid := domain.ProjectId(uuid.MustParse(projectId))
	menuItemUuid := domain.MenuItemId(uuid.MustParse(menuItemId))

	tx, err := s.characterRepository.BeginTransaction(context.TODO())

	if err != nil {
		return err
	}
	defer tx.Rollback()

	character, err := s.characterRepository.FindCharacter(tx, userUUid, projectUUid, menuItemUuid)
	if err != nil {
		return err
	}

	character.Name = name
	if description != nil {
		character.Description = description
	}

	var charAttrs []domain.Attribute
	for _, attribute := range attributes {
		charAttrs = append(charAttrs, domain.Attribute{
			Field: attribute.Field,
			Value: attribute.Value,
		})
	}
	character.Attributes = charAttrs

	err = s.characterRepository.Save(tx, character)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateLocation(
	userId string,
	projectId string,
	menuItemId string,
	name string,
	description *string,
) error {
	userUUid := domain.UserId(uuid.MustParse(userId))
	projectUUid := domain.ProjectId(uuid.MustParse(projectId))
	menuItemUuid := domain.MenuItemId(uuid.MustParse(menuItemId))

	tx, err := s.locationRepository.BeginTransaction(context.TODO())

	if err != nil {
		return err
	}
	defer tx.Rollback()

	location, err := s.locationRepository.FindLocation(tx, userUUid, projectUUid, menuItemUuid)
	if err != nil {
		return err
	}

	location.Name = name
	if description != nil {
		location.Description = description
	}

	err = s.locationRepository.Save(tx, location)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProject(userId string, projectId string) (readmodel.Project, error) {
	return s.projectReadModel.GetById(userId, projectId)
}

func (s *Service) UpdateProject(userId string, projectId string, title string, description *string) error {
	userUUid := domain.UserId(uuid.MustParse(userId))
	projectUUid := domain.ProjectId(uuid.MustParse(projectId))

	tx, err := s.projectRepository.BeginTransaction(context.TODO())
	if err != nil {
		return err
	}
	defer tx.Rollback()

	exists, err := s.projectRepository.ExistsByUserIdAndTitle(tx, userUUid, title, &projectUUid)

	if err != nil {
		return err
	}

	if exists {
		return gerr.NewBusinessLogicErr("Project with this title already exists", gerr.ProjectAlreadyExists)
	}

	project, err := s.projectRepository.FindByUserIdAndProjectId(tx, userUUid, projectUUid)

	if err != nil {
		return err
	}

	if project == nil {
		return gerr.NewBusinessLogicErr("project not found", gerr.NotFound)
	}

	project.Title = title
	if description != nil {
		project.Description = *description
	} else {
		project.Description = ""
	}

	err = s.projectRepository.Save(tx, project)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
