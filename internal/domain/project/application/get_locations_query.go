package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type GetLocationsQuery struct {
	UserId    string
	ProjectId string
}

type GetLocationsHandlerResponse struct {
	Locations []readmodel.Location
}

type GetLocationsHandler struct {
	locationReadModel readmodel.LocationReadModel
}

func NewGetLocationsHandler(locationReadModel readmodel.LocationReadModel) *GetLocationsHandler {
	return &GetLocationsHandler{locationReadModel: locationReadModel}
}

func (h *GetLocationsHandler) Handle(ctx context.Context, query *GetLocationsQuery) (*GetLocationsHandlerResponse, error) {
	locations, err := h.locationReadModel.GetLocations(query.UserId, query.ProjectId)

	if err != nil {
		return nil, err
	}

	return &GetLocationsHandlerResponse{Locations: locations}, nil
}
