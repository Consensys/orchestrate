package envelope

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/opentracing/opentracing-go"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
)

// Error returns string representation of errors encountered by envelope
func (e *Envelope) Error() string {
	return common.Errors(e.Errors).Error()
}

// Sender returns sender of the transaction
func (e *Envelope) Sender() ethcommon.Address {
	return e.GetFrom().Address()
}

// Carrier returns an OpenTracing carrier based on envelope Metadata
func (e *Envelope) Carrier() opentracing.TextMapCarrier {
	if e.GetMetadata() == nil {
		e.Metadata = &Metadata{
			Extra: make(map[string]string),
		}
	} else if e.GetMetadata().GetExtra() == nil {
		e.GetMetadata().Extra = make(map[string]string)
	}

	return opentracing.TextMapCarrier(e.GetMetadata().GetExtra())
}
