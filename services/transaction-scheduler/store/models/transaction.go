package models

import "time"

type Transaction struct {
	tableName struct{} `pg:"transactions"` // nolint:unused,structcheck // reason

	ID        int
	UUID      string
	From      string
	To        *string
	Hash      string
	Type      string
	Data      *string
	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time
}
