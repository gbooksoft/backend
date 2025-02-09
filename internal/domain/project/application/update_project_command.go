package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type UpdateProjectCommand struct {
	UserId      string
	ProjectId   string
	Title       string
	Description *string
	Genres      []string
	Tags        []string
}

type UpdateProjectHandlerResponse struct {
}

type UpdateProjectHandler struct {
	projectRepository domain.ProjectRepository
	genreRepository   domain.GenreRepository
}

func NewUpdateProjectHandler(projectRepository domain.ProjectRepository, genreRepository domain.GenreRepository) *UpdateProjectHandler {
	return &UpdateProjectHandler{
		projectRepository: projectRepository,
		genreRepository:   genreRepository,
	}
}

func (h *UpdateProjectHandler) Handle(ctx context.Context, cmd *UpdateProjectCommand) (*UpdateProjectHandlerResponse, error) {
	userUUid := domain.UserId(uuid.MustParse(cmd.UserId))
	projectUUid := domain.ProjectId(uuid.MustParse(cmd.ProjectId))
	title := cmd.Title
	description := cmd.Description

	tx, err := h.projectRepository.BeginTransaction(context.TODO())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	exists, err := h.projectRepository.ExistsByUserIdAndTitle(tx, userUUid, title, &projectUUid)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, gerr.NewBusinessLogicErr("Project with this title already exists", gerr.ProjectAlreadyExists)
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

	project, err := h.projectRepository.FindByUserIdAndProjectId(tx, userUUid, projectUUid)

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, gerr.NewBusinessLogicErr("project not found", gerr.NotFound)
	}

	project.Title = title
	if description != nil {
		project.Description = *description
	} else {
		project.Description = ""
	}

	project.GenresIds = genres
	project.Tags = cmd.Tags

	err = h.projectRepository.Save(tx, project)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return nil, nil
}
