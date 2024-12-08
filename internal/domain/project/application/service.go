package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
	"github.com/gbooksoft/backend/internal/gerr"
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
	}
}

func (s *Service) CreateProject(userId string, title, description string) (string, string, error) {
	tx, err := s.projectRepository.BeginTransaction(context.TODO())

	if err != nil {
		return "", "", err
	}
	defer tx.Rollback()

	userUuid := domain.UserId(uuid.MustParse(userId))

	exists, err := s.projectRepository.ExistsByUserIdAndTitle(tx, userUuid, title)

	if err != nil {
		return "", "", err
	}

	if exists {
		return "", "", gerr.NewBusinessLogicErr("Project with this title already exists", gerr.ProjectAlreadyExists, nil)
	}

	projectId := domain.ProjectId(s.projectRepository.GetNextId())

	project := domain.NewProject(projectId, userUuid, title, description)

	err = s.projectRepository.Add(tx, project)

	if err != nil {
		return "", "", err
	}

	menu, expandedItem := domain.NewDefaultMenu(projectId)

	err = s.menuRepository.AddMenu(tx, menu)

	if err != nil {
		return "", "", err
	}

	editors := domain.NewEditorsFromMenu(menu)

	err = s.editorRepository.AddMany(tx, editors)

	if err != nil {
		return "", "", err
	}

	if err = tx.Commit(); err != nil {
		return "", "", err
	}

	return projectId.String(), expandedItem.String(), nil
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

func (s *Service) AddMenuItem(newMenuItemId, userId, projectId, parentId string) error {
	projectUUid := domain.ProjectId(uuid.MustParse(projectId))
	newMenuItemUUid := domain.MenuItemId(uuid.MustParse(newMenuItemId))
	parentUUid := domain.MenuItemId(uuid.MustParse(parentId))
	userUUid := domain.UserId(uuid.MustParse(userId))

	tx, err := s.menuRepository.BeginTransaction(context.TODO())

	if err != nil {
		return err
	}
	defer tx.Rollback()

	menu, err := s.menuRepository.FindByProjectId(tx, userUUid, projectUUid)

	if err != nil {
		return err
	}

	newMenuItem, err := menu.AddNewItem(newMenuItemUUid, parentUUid)

	if err != nil {
		return err
	}

	err = s.menuRepository.Save(tx, userUUid, projectUUid, menu)

	if err != nil {
		return err
	}

	if newMenuItem.Type == domain.TypeChapter {
		editor := domain.NewEditor(newMenuItemUUid)

		err = s.editorRepository.AddEditor(tx, editor)

		if err != nil {
			return err
		}
	}

	if newMenuItem.Type == domain.TypeCharacter {
		character := domain.NewCharacter(newMenuItemUUid)

		err = s.characterRepository.AddCharacter(tx, character)

		if err != nil {
			return err
		}
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

func (s *Service) UpdateMenuItem(userId string, projectId string, menuItemId string, title string, setSelected bool) error {
	userUUid := domain.UserId(uuid.MustParse(userId))
	projectUUid := domain.ProjectId(uuid.MustParse(projectId))
	menuItemUuid := domain.MenuItemId(uuid.MustParse(menuItemId))

	tx, err := s.menuRepository.BeginTransaction(context.TODO())
	if err != nil {
		return err
	}
	defer tx.Rollback()

	menu, err := s.menuRepository.FindByProjectId(tx, userUUid, projectUUid)
	if err != nil {
		return err
	}

	if len(title) > 0 {
		err = menu.UpdateItemTitle(menuItemUuid, title)

		if err != nil {
			return err
		}
	}

	if setSelected {
		menu.SelectItem(menuItemUuid)
	}

	err = s.menuRepository.Save(tx, userUUid, projectUUid, menu)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) SaveEditorState(userId, projectId, menuItemId, editorState string) error {
	userUUid := domain.UserId(uuid.MustParse(userId))
	projectUUid := domain.ProjectId(uuid.MustParse(projectId))
	menuItemUuid := domain.MenuItemId(uuid.MustParse(menuItemId))

	tx, err := s.editorRepository.BeginTransaction(context.TODO())

	if err != nil {
		return err
	}
	defer tx.Rollback()

	editor, err := s.editorRepository.FindEditor(tx, userUUid, projectUUid, menuItemUuid)

	if err != nil {
		return err
	}

	editor.EditorState = &editorState

	err = s.editorRepository.Save(tx, editor)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) FetchEditors(userId string, projectId string) ([]readmodel.Editor, error) {
	return s.editorReadModel.GetEditors(userId, projectId)
}

func (s *Service) FetchCharacters(userId string, projectId string) ([]readmodel.Character, error) {
	return s.characterReadModel.GetCharacters(userId, projectId)
}

func (s *Service) GetAllProjects(userId string) ([]readmodel.Project, error) {
	return s.projectReadModel.GetAll(userId)
}

func (s *Service) GetRecentProjects(userId string) ([]readmodel.Project, error) {
	return s.projectReadModel.GetRecent(userId, 3)
}

func (s *Service) DeleteMenuItem(userId string, projectId string, menuItemId string) error {
	userUUid := domain.UserId(uuid.MustParse(userId))
	projectUUid := domain.ProjectId(uuid.MustParse(projectId))
	menuItemUuid := domain.MenuItemId(uuid.MustParse(menuItemId))

	err := s.menuRepository.DeleteMenuItem(userUUid, projectUUid, menuItemUuid)

	if err != nil {
		return err
	}

	return nil
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
	description string,
	attributes []CharacterAttribute,
) error {
	userUUid := domain.UserId(uuid.MustParse(userId))
	projectUUid := domain.ProjectId(uuid.MustParse(projectId))
	menuItemUuid := domain.MenuItemId(uuid.MustParse(menuItemId))

	tx, err := s.editorRepository.BeginTransaction(context.TODO())

	if err != nil {
		return err
	}
	defer tx.Rollback()

	character, err := s.characterRepository.FindCharacter(tx, userUUid, projectUUid, menuItemUuid)
	if err != nil {
		return err
	}

	character.Name = name
	character.Description = description

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

func (s *Service) GetProject(userId string, projectId string) (readmodel.Project, error) {
	return s.projectReadModel.GetById(userId, projectId)
}
