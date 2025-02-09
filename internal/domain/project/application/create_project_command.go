package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type CreateProjectCommand struct {
	UserId       string
	UserFullName string
	Title        string
	Description  string
	Genres       []string
	Tags         []string
}

type CreateProjectHandlerResponse struct {
	ProjectId      string
	ItemIdToExpand string
}

type CreateProjectHandler struct {
	projectRepository domain.ProjectRepository
	menuRepository    domain.MenuRepository
	editorRepository  domain.EditorRepository
	genreRepository   domain.GenreRepository
}

func NewCreateProjectHandler(
	projectRepository domain.ProjectRepository,
	menuRepository domain.MenuRepository,
	editorRepository domain.EditorRepository,
	genreRepository domain.GenreRepository,
) *CreateProjectHandler {
	return &CreateProjectHandler{
		projectRepository: projectRepository,
		menuRepository:    menuRepository,
		editorRepository:  editorRepository,
		genreRepository:   genreRepository,
	}
}

func (h *CreateProjectHandler) Handle(ctx context.Context, cmd *CreateProjectCommand) (*CreateProjectHandlerResponse, error) {
	title := cmd.Title
	description := cmd.Description
	userId := cmd.UserId
	userFullName := cmd.UserFullName
	tags := cmd.Tags

	tx, err := h.projectRepository.BeginTransaction(context.TODO())

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userUuid := domain.UserId(uuid.MustParse(userId))

	exists, err := h.projectRepository.ExistsByUserIdAndTitle(tx, userUuid, title, nil)

	if err != nil {
		return nil, err
	}

	genresExists, err := h.genreRepository.Exists(tx, cmd.Genres)
	if err != nil {
		return nil, err
	}

	if !genresExists {
		return nil, gerr.NewBusinessLogicErr("Not all of given genres are exists", gerr.GenreNotExists)
	}

	var genres []domain.GenreId
	for _, genre := range cmd.Genres {
		genres = append(genres, domain.GenreId(uuid.MustParse(genre)))
	}

	if exists {
		return nil, gerr.NewBusinessLogicErr("Project with this title already exists", gerr.ProjectAlreadyExists)
	}

	projectId := domain.ProjectId(h.projectRepository.GetNextId())

	project := domain.NewProject(projectId, userUuid, title, description, genres, tags)

	err = h.projectRepository.Add(tx, project)

	if err != nil {
		return nil, err
	}

	metadataId := domain.MetadataId(h.projectRepository.GetNextId())
	metadata := domain.NewMetadata(metadataId, projectId, userFullName, "", "", "", "")
	err = h.projectRepository.AddMetadata(tx, metadata)

	if err != nil {
		return nil, err
	}

	menu, expandedItem := domain.NewDefaultMenu(projectId)

	err = h.menuRepository.AddMenu(tx, menu)

	if err != nil {
		return nil, err
	}

	editors := domain.NewEditorsFromMenu(menu)

	err = h.editorRepository.AddMany(tx, editors)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &CreateProjectHandlerResponse{
		ProjectId:      projectId.String(),
		ItemIdToExpand: expandedItem.String(),
	}, nil
}
