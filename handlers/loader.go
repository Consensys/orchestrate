package handlers

import (
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/services"
)

// Loader creates an handler loading input
func Loader(u services.Unmarshaller) types.HandlerFunc {
	pool := &sync.Pool{
		New: func() interface{} { return &tracepb.Trace{} },
	}

	return func(ctx *types.Context) {
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
