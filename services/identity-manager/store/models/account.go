package models

import (
	"time"
)

type Account struct {
	tableName struct{} `pg:"accounts"` // nolint:unused,structcheck // reason

	ID                  int
	Alias               string
	Address             string
	PublicKey           string `pg:"alias:public_key"`
	CompressedPublicKey string `pg:"alias:compressed_public_key"`
	TenantID            string `pg:"alias:tenant_id"`
	Attributes          map[string]string

	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time `pg:"default:now()"`
}
