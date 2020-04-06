package models

import (
	"fmt"
	"strings"
	"time"

	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

// EnvelopeModel represent elements in `envelopes` table
type EnvelopeModel struct {
	tableName struct{} `pg:"envelopes"` //nolint:unused,structcheck // reason

	// ID technical identifier
	ID int32

	// Envelope Identifier
	EnvelopeID string

	// Tenant Identifier
	TenantID string

	// Ethereum info about transaction
	ChainID string
	TxHash  string

	// Envelope
	Envelope []byte

	// Status
	Status   string
	StoredAt time.Time
	SentAt   time.Time
	MinedAt  time.Time
	ErrorAt  time.Time
}

// FromEnvelope creates a model from an envelope
func NewEnvelopeFromTx(tenantID string, e *tx.TxEnvelope) (EnvelopeModel, error) {
	// Marshal envelope
	b, err := encoding.Marshal(e)
	if err != nil {
		return EnvelopeModel{}, err
	}

	return EnvelopeModel{
		Envelope:   b,
		EnvelopeID: e.GetID(),
		TenantID:   tenantID,
		ChainID:    e.GetChainID(),
		TxHash:     e.GetTxHash(),
	}, nil
}

func NewEnvelope(tenantId string, chainId string, txHash string) (EnvelopeModel) {
	return EnvelopeModel{
		ChainID:  chainId,
		TenantID: tenantId,
		TxHash:   txHash,
	}
}

// StatusInfo returns a proto formatted StatusInfo object from a model
func (model *EnvelopeModel) StatusInfo() *svc.StatusInfo {
	return &svc.StatusInfo{
		Status:   model.StatusFormatted(),
		StoredAt: utils.TimeToPTimestamp(model.StoredAt),
		SentAt:   utils.TimeToPTimestamp(model.SentAt),
		MinedAt:  utils.TimeToPTimestamp(model.MinedAt),
		ErrorAt:  utils.TimeToPTimestamp(model.ErrorAt),
	}
}

// StatusFormatted returns a proto Status enum
func (model *EnvelopeModel) StatusFormatted() svc.Status {
	status, ok := svc.Status_value[strings.ToUpper(model.Status)]
	if !ok {
		panic(fmt.Sprintf("invalid status %q", model.Status))
	}
	return svc.Status(status)
}