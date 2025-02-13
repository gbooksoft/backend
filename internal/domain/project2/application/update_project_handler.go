package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project2/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type UpdateProjectCommand struct {
	AuthorId    string
	ProjectId   string
	Title       string
	Description *string
	Genres      []string
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
	ctx, err := h.projectRepository.BeginTransaction(context.TODO())
	if err != nil {
		return nil, err
	}
	defer h.projectRepository.RollbackTransaction(ctx)

	authorId := domain.AuthorId(uuid.MustParse(cmd.AuthorId))
	projectId := domain.ProjectId(uuid.MustParse(cmd.ProjectId))

	exists, err := h.projectRepository.ExistsByUserIdAndTitle(ctx, authorId, cmd.Title, &projectId)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, gerr.NewBusinessLogicErr("Project with this title already exists", gerr.ProjectAlreadyExists)
	}

	var genres []domain.Genre
	for _, genre := range cmd.Genres {
		genres = append(genres, domain.Genre{Id: domain.GenreId(uuid.MustParse(genre))})
	}

	project, err := h.projectRepository.FindByAuthorIdAndProjectId(ctx, authorId, projectId)

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, gerr.NewBusinessLogicErr("project not found", gerr.NotFound)
	}

	err = project.Update(cmd.Title, cmd.Description, genres)

	if err != nil {
		return nil, err
	}

	err = h.projectRepository.Save(ctx, project)

	if err != nil {
		return nil, err
	}

	if err := h.projectRepository.CommitTransaction(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
