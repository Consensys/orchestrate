package api

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

type TransactionResponse struct {
	UUID           string                         `json:"uuid" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	IdempotencyKey string                         `json:"idempotencyKey" example:"myIdempotencyKey"`
	ChainName      string                         `json:"chain" example:"myChain"`
	Params         *entities.ETHTransactionParams `json:"params"`
	Jobs           []*JobResponse                 `json:"jobs"`
	CreatedAt      time.Time                      `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
}
