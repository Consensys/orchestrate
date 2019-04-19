package envelope

import (
	"github.com/opentracing/opentracing-go"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
)

// Error returns string representation of errors encountered by envelope
func (e *Envelope) Error() string {
	return common.Errors(e.Errors).Error()
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
