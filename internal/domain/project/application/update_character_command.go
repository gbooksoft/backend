package application

import (
	"context"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
	"github.com/google/uuid"
)

type CharacterAttribute struct {
	Field string
	Value string
}

type UpdateCharacterCommand struct {
	UserId      string
	ProjectId   string
	MenuItemId  string
	Name        string
	Description *string
	Attributes  []CharacterAttribute
}

type UpdateCharacterHandlerResponse struct {
}

type UpdateCharacterHandler struct {
	characterRepository domain.CharacterRepository
}

func NewUpdateCharacterHandler(characterRepository domain.CharacterRepository) *UpdateCharacterHandler {
	return &UpdateCharacterHandler{
		characterRepository: characterRepository,
	}
}

func (h *UpdateCharacterHandler) Handle(ctx context.Context, cmd *UpdateCharacterCommand) (*UpdateCharacterHandlerResponse, error) {
	userUUid := domain.UserId(uuid.MustParse(cmd.UserId))
	projectUUid := domain.ProjectId(uuid.MustParse(cmd.ProjectId))
	menuItemUuid := domain.MenuItemId(uuid.MustParse(cmd.MenuItemId))
	name := cmd.Name
	description := cmd.Description
	attributes := cmd.Attributes

	ctx, err := h.characterRepository.BeginTransaction(context.TODO())

	if err != nil {
		return nil, err
	}
	tx, ok := repository.FetchDBTransactionFromContext(ctx)
	if !ok {
		panic("transaction not found")
	}
	defer tx.Rollback()

	character, err := h.characterRepository.FindCharacter(tx, userUUid, projectUUid, menuItemUuid)
	if err != nil {
		return nil, err
	}

	character.Name = name
	if description != nil {
		character.Description = description
	}

	var charAttrs []domain.Attribute
	for _, attribute := range attributes {
		charAttrs = append(charAttrs, domain.Attribute{
			Field: attribute.Field,
			Value: attribute.Value,
		})
	}
	character.Attributes = charAttrs

	err = h.characterRepository.Save(tx, character)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
