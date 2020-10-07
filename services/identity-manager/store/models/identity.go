package models

import (
	"time"
)

type Identity struct {
	tableName struct{} `pg:"identities"` // nolint:unused,structcheck // reason

	ID        int
	UUID      string
	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time `pg:"default:now()"`
}
