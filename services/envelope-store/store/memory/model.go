package memory

import (
	"fmt"
	"strings"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

// EnvelopeModel is a entry into mock envelope Store
type EnvelopeModel struct {
	// envelope
	envelope *tx.TxEnvelope

	// Status
	Status   string
	StoredAt time.Time
	ErrorAt  time.Time
	SentAt   time.Time
	MinedAt  time.Time
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

func (model *EnvelopeModel) ToStatusResponse() (*svc.StatusResponse, error) {
	return &svc.StatusResponse{
		StatusInfo: model.StatusInfo(),
	}, nil
}

func (model *EnvelopeModel) ToStoreResponse() (*svc.StoreResponse, error) {
	return &svc.StoreResponse{
		StatusInfo: model.StatusInfo(),
		Envelope:   model.envelope,
	}, nil

}
