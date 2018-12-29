package handlers

import (
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

// Producer creates a producer handler
func Producer(p infra.TraceProducer) infra.HandlerFunc {

	pool := &sync.Pool{
		New: func() interface{} { return &tracepb.Trace{} },
	}

	return func(ctx *infra.Context) {
		pb := pool.Get().(*tracepb.Trace)
		defer pool.Put(pb)

		pb.Reset()
		protobuf.DumpTrace(ctx.T, pb)
		
		err := p.Produce(pb)
		if err != nil {
			ctx.Error(err)
		}
	}
}
