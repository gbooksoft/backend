package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
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
}

func NewCreateProjectHandler(projectRepository domain.ProjectRepository) *CreateProjectHandler {
	return &CreateProjectHandler{
		projectRepository: projectRepository,
	}
}

func (h *CreateProjectHandler) Handle(ctx context.Context, cmd *CreateProjectCommand) (*CreateProjectHandlerResponse, error) {
	ctx, err := h.projectRepository.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer h.projectRepository.RollbackTransaction(ctx)

	authorId := domain.AuthorId(uuid.MustParse(cmd.AuthorId))

	exists, err := h.projectRepository.ExistsByUserIdAndTitle(ctx, authorId, cmd.Title, nil)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, gerr.NewBusinessLogicErr("project already exists", gerr.ProjectAlreadyExists)
	}

	var genres []domain.Genre
	for _, genre := range cmd.Genres {
		genres = append(genres, domain.Genre{Id: domain.GenreId(uuid.MustParse(genre))})
	}

	projectId := domain.ProjectId(h.projectRepository.GetNextId())

	project, err := domain.NewProject(projectId, cmd.Title, cmd.Description, genres, authorId)

	if err != nil {
		return nil, err
	}

	err = h.projectRepository.Add(ctx, project)

	if err != nil {
		return nil, err
	}

	if err := h.projectRepository.CommitTransaction(ctx); err != nil {
		return nil, err
	}

	return &CreateProjectHandlerResponse{
		ProjectId: projectId.String(),
	}, nil
}
