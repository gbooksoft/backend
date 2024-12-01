package domain

import (
	"database/sql"
	"github.com/gbooksoft/backend/pkg/repository"
)

type CharacterRepository interface {
	repository.BaseRepository
	AddCharacter(tx *sql.Tx, character *Character) error
}
