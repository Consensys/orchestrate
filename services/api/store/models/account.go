package models

import (
	"time"
)

type Account struct {
	tableName struct{} `pg:"accounts"` // nolint:unused,structcheck // reason

	ID                  int
	Alias               string
	Address             string
	PublicKey           string
	CompressedPublicKey string
	TenantID            string
	Attributes          map[string]string

	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time `pg:"default:now()"`
}
