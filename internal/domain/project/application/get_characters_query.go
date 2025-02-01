package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type GetCharactersQuery struct {
	UserId    string
	ProjectId string
}

type GetCharactersHandlerResponse struct {
	Characters []readmodel.Character
}

type GetCharactersHandler struct {
	characterReadModel readmodel.CharacterReadModel
}

func NewGetCharactersHandler(characterReadModel readmodel.CharacterReadModel) *GetCharactersHandler {
	return &GetCharactersHandler{characterReadModel: characterReadModel}
}

func (h *GetCharactersHandler) Handle(ctx context.Context, query *GetCharactersQuery) (*GetCharactersHandlerResponse, error) {
	characters, err := h.characterReadModel.GetCharacters(query.UserId, query.ProjectId)

	if err != nil {
		return nil, err
	}

	return &GetCharactersHandlerResponse{Characters: characters}, nil
}
