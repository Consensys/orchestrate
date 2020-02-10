//nolint:stylecheck // reason
package envelope

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"
)

func (e *Envelope) TxRequest() *tx.TxRequest {
	return e.Body.(*Envelope_TxRequest).TxRequest
}

func (e *Envelope) TxResponse() *tx.TxResponse {
	return e.Body.(*Envelope_TxResponse).TxResponse
}

func (e *Envelope) Builder() (*tx.Builder, error) {
	switch e.Body.(type) {
	case *Envelope_TxRequest:
		return e.TxRequest().Builder()
	case *Envelope_TxResponse:
		return e.TxResponse().Builder()
	default:
		return nil, errors.DataError("")
	}
}
