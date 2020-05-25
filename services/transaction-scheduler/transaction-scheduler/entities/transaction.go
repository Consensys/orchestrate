package entities

import (
	"time"
)

type Transaction struct {
	ID             int
	Hash           string
	From           string
	To             string
	Nonce          string
	Value          string
	GasPrice       string
	GasLimit       string
	Data           string
	Raw            string
	PrivateFrom    string
	PrivateFor     []string
	PrivacyGroupID string
	CreatedAt      time.Time
}
