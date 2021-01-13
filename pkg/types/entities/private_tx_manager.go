package entities

import "time"

type PrivateTxManager struct {
	UUID      string
	ChainUUID string
	URL       string
	Type      string
	CreatedAt time.Time
}
