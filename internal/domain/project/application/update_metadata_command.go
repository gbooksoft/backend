package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/google/uuid"
)

type UpdateMetadataCommand struct {
	UserId    string
	ProjectId string
	Author    string
	Genre     string
	Keywords  string
	Language  string
	Isbn      string
}

type UpdateMetadataHandlerResponse struct {
}

type UpdateMetadataHandler struct {
	projectRepository domain.ProjectRepository
}

func NewUpdateMetadataHandler(projectRepository domain.ProjectRepository) *UpdateMetadataHandler {
	return &UpdateMetadataHandler{
		projectRepository: projectRepository,
	}
}

func (h *UpdateMetadataHandler) Handle(ctx context.Context, cmd *UpdateMetadataCommand) (*UpdateMetadataHandlerResponse, error) {
	userUUid := domain.UserId(uuid.MustParse(cmd.UserId))
	projectUUid := domain.ProjectId(uuid.MustParse(cmd.ProjectId))

	tx, err := h.projectRepository.BeginTransaction(context.TODO())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	metadata, err := h.projectRepository.FindMetadataByUserIdAndProjectId(tx, userUUid, projectUUid)

	if err != nil {
		return nil, err
	}

	if metadata == nil {
		return nil, gerr.NewBusinessLogicErr("metadata not found", gerr.ProjectMetadataNotFound)
	}

	if exists, err := h.projectRepository.ExistsByUserIdAndProjectId(tx, userUUid, metadata.ProjectId); err != nil || !exists {
		if err != nil {
			return nil, err
		}

		return nil, gerr.NewBusinessLogicErr("project not found", gerr.NotFound)
	}

	metadata.Author = cmd.Author
	metadata.Genre = cmd.Genre
	metadata.Keywords = cmd.Keywords
	metadata.Language = cmd.Language
	metadata.Isbn = cmd.Isbn

	err = h.projectRepository.SaveMetadata(tx, metadata)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return nil, nil
}
