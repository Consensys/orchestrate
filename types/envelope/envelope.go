package envelope

import (
	"github.com/opentracing/opentracing-go"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
)

func (t *Envelope) Error() string {
	return common.Errors(t.Errors).Error()
}

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