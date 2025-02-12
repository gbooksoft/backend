package model

import "github.com/google/uuid"

type Uuid uuid.UUID

func (id Uuid) String() string {
	return uuid.UUID(id).String()
}

func (id Uuid) Equal(u Uuid) bool {
	return uuid.UUID(id) == uuid.UUID(u)
}
