package models

import (
	"time"
)

type Transaction struct {
	tableName struct{} `pg:"transactions"` // nolint:unused,structcheck // reason

	ID             int
	UUID           string
	Hash           string
	Sender         string
	Recipient      string
	Nonce          string
	Value          string
	GasPrice       string
	Gas            string
	Data           string
	Raw            string
	PrivateFrom    string
	PrivateFor     []string `pg:",array"`
	PrivacyGroupID string
	CreatedAt      time.Time `pg:"default:now()"`
	UpdatedAt      time.Time `pg:"default:now()"`
}
