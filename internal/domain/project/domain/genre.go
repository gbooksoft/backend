package domain

import "github.com/google/uuid"

type GenreId uuid.UUID

func (id GenreId) String() string {
	return uuid.UUID(id).String()
}
