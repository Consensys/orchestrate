package mock

import (
	"fmt"
	"strings"
	"time"

	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
)

// EnvelopeModel is a entry into mock envelope Store
type EnvelopeModel struct {
	// envelope
	envelope *envelope.Envelope

	// Status
	Status   string
	StoredAt time.Time
	ErrorAt  time.Time
	SentAt   time.Time
	MinedAt  time.Time
}

// StatusInfo returns a proto formated StatusInfo object from a model
func (model *EnvelopeModel) StatusInfo() *evlpstore.StatusInfo {
	return &evlpstore.StatusInfo{
		Status:   model.StatusFormated(),
		StoredAt: utils.TimeToPTimestamp(model.StoredAt),
		SentAt:   utils.TimeToPTimestamp(model.SentAt),
		MinedAt:  utils.TimeToPTimestamp(model.MinedAt),
		ErrorAt:  utils.TimeToPTimestamp(model.ErrorAt),
	}
}

// StatusFormated returns a proto Status enum
func (model *EnvelopeModel) StatusFormated() evlpstore.Status {
	status, ok := evlpstore.Status_value[strings.ToUpper(model.Status)]
	if !ok {
		panic(fmt.Sprintf("invalid status %q", model.Status))
	}
	return evlpstore.Status(status)
}

func (model *EnvelopeModel) ToStatusResponse() (*evlpstore.StatusResponse, error) {
	return &evlpstore.StatusResponse{
		StatusInfo: model.StatusInfo(),
	}, nil
}

func (model *EnvelopeModel) ToStoreResponse() (*evlpstore.StoreResponse, error) {
	return &evlpstore.StoreResponse{
		StatusInfo: model.StatusInfo(),
		Envelope:   model.envelope,
	}, nil

}
