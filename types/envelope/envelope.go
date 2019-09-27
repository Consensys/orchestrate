//nolint:stylecheck
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

// maybeInitMetadata initialize metadata object with Extra field
func (e *Envelope) maybeInitMetadata() {
	if e.GetMetadata() == nil {
		e.Metadata = &Metadata{
			Extra: make(map[string]string),
		}
	} else if e.GetMetadata().GetExtra() == nil {
		e.GetMetadata().Extra = make(map[string]string)
	}
}

// Carrier returns an OpenTracing carrier based on envelope Metadata
func (e *Envelope) Carrier() opentracing.TextMapCarrier {
	e.maybeInitMetadata()
	return opentracing.TextMapCarrier(e.GetMetadata().GetExtra())
}

// GetMetadataValue retrieves value stored in Metadata extra
func (e *Envelope) GetMetadataValue(key string) (string, bool) {
	if e.GetMetadata() == nil || e.GetMetadata().GetExtra() == nil {
		return "", false
	}
	v, ok := e.GetMetadata().GetExtra()[key]
	return v, ok
}

// SetMetadataValue set a value stored in Metadata extra
func (e *Envelope) SetMetadataValue(key, value string) {
	e.maybeInitMetadata()
	e.GetMetadata().GetExtra()[key] = value
}
