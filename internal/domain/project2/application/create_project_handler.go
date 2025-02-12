package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project2/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type CreateProjectCommand struct {
	AuthorId     string
	UserFullName string
	Title        string
	Description  string
	Genres       []string
	Tags         []string
}

type CreateProjectHandlerResponse struct {
	ProjectId string
}

type CreateProjectHandler struct {
	projectRepository domain.ProjectRepository
	genreRepository   domain.GenreRepository
}

func NewCreateProjectHandler(projectRepository domain.ProjectRepository, genreRepository domain.GenreRepository) *CreateProjectHandler {
	return &CreateProjectHandler{
		projectRepository: projectRepository,
		genreRepository:   genreRepository,
	}
}

func (h *CreateProjectHandler) Handle(ctx context.Context, cmd *CreateProjectCommand) (*CreateProjectHandlerResponse, error) {
	tx, err := h.projectRepository.BeginTransaction(ctx)

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	authorId := domain.AuthorId(uuid.MustParse(cmd.AuthorId))

	exists, err := h.projectRepository.ExistsByUserIdAndTitle(tx, authorId, cmd.Title)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, gerr.NewBusinessLogicErr("project already exists", gerr.ProjectAlreadyExists)
	}

	var genres []domain.GenreId
	for _, genre := range cmd.Genres {
		genres = append(genres, domain.GenreId(uuid.MustParse(genre)))
	}

	genresExists, err := h.genreRepository.Exists(tx, genres)

	if err != nil {
		return nil, err
	}

	if !genresExists {
		return nil, gerr.NewBusinessLogicErr("Not all of given genres are exists", gerr.GenreNotExists)
	}

	projectId := domain.ProjectId(h.projectRepository.GetNextId())

	project, err := domain.NewProject(projectId, cmd.Title, cmd.Description, genres, authorId)

	if err != nil {
		return nil, err
	}

	err = h.projectRepository.Add(tx, project)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &CreateProjectHandlerResponse{
		ProjectId: projectId.String(),
	}, nil
}
