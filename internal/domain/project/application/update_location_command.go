package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/google/uuid"
)

type UpdateLocationCommand struct {
	UserId      string
	ProjectId   string
	MenuItemId  string
	Name        string
	Description *string
}

type UpdateLocationHandlerResponse struct {
}

type UpdateLocationHandler struct {
	locationRepository domain.LocationRepository
}

func NewUpdateLocationHandler(locationRepository domain.LocationRepository) *UpdateLocationHandler {
	return &UpdateLocationHandler{
		locationRepository: locationRepository,
	}
}

func (h *UpdateLocationHandler) Handle(ctx context.Context, cmd *UpdateLocationCommand) (*UpdateLocationHandlerResponse, error) {
	userUUid := domain.UserId(uuid.MustParse(cmd.UserId))
	projectUUid := domain.ProjectId(uuid.MustParse(cmd.ProjectId))
	menuItemUuid := domain.MenuItemId(uuid.MustParse(cmd.MenuItemId))
	name := cmd.Name
	description := cmd.Description

	tx, err := h.locationRepository.BeginTransaction(context.TODO())

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	location, err := h.locationRepository.FindLocation(tx, userUUid, projectUUid, menuItemUuid)
	if err != nil {
		return nil, err
	}

	location.Name = name
	if description != nil {
		location.Description = description
	}

	err = h.locationRepository.Save(tx, location)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
