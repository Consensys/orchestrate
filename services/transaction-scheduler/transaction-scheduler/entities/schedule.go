package entities

import (
	"time"
)

type Schedule struct {
	UUID      string
	ChainUUID string
	Jobs      []*Job
	TxRequest *TxRequest
	CreatedAt time.Time
}
