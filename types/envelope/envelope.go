package envelope

import (
	fmt "fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/opentracing/opentracing-go"
)

// Error returns string representation of errors encountered by envelope
func (e *Envelope) Error() string {
	if len(e.GetErrors()) == 0 {
		return ""
	}
	return fmt.Sprintf("%q", e.GetErrors())
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
