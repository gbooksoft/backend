package domain

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type CharacterRepository interface {
	repository.BaseRepository
	AddCharacter(tx *sql.Tx, character *Character) error
	FindCharacter(tx *sql.Tx, userId UserId, projectId ProjectId, menuItemId MenuItemId) (*Character, error)
	Save(tx *sql.Tx, character *Character) error
}
