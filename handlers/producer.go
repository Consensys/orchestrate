package handlers

import (
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Producer creates a producer handler
func Producer(p services.TraceProducer) types.HandlerFunc {

	pool := &sync.Pool{
		New: func() interface{} { return &tracepb.Trace{} },
	}

	return func(ctx *types.Context) {
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
