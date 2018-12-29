package handlers

import (
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

// Loader creates an handler loading input
func Loader(u infra.Unmarshaller) infra.HandlerFunc {
	pool := &sync.Pool{
		New: func() interface{} { return &tracepb.Trace{} },
	}

	return func(ctx *infra.Context) {
		pb := pool.Get().(*tracepb.Trace)
		defer pool.Put(pb)

		pb.Reset()

		// Unmarshal message
		err := u.Unmarshal(ctx.Msg, pb)
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}

		// Load Trace from protobuffer
		protobuf.LoadTrace(pb, ctx.T)
	}
}
