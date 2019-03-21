package handlers

import (
	"context"
	"math/big"
	"math/rand"
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

var letterRunes = []rune("abcdef0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func makeTraceLoaderContext(i int, store *mock.TraceStore) *worker.Context {
	ctx := worker.NewContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())

	hash, uuid := "0x"+RandString(32), RandString(32)

	// Prepare context
	ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(10))
	ctx.T.Receipt = &ethereum.Receipt{
		TxHash: hash,
	}

	switch i % 2 {
	case 0:
		// Prestore a trace
		tr := &trace.Trace{}
		tr.Chain = (&common.Chain{}).SetID(big.NewInt(10))
		tr.Tx = (&ethereum.Transaction{
			Hash: hash,
			Raw:  "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
		})
		tr.Metadata = &trace.Metadata{Id: uuid}
		store.Store(context.Background(), tr)
		ctx.Keys["uuid"] = uuid
		ctx.Keys["errors"] = 0
	case 1:
		ctx.Keys["uuid"] = ""
		ctx.Keys["errors"] = 1
	}

	return ctx
}

func TestTraceLoader(t *testing.T) {
	store := mock.NewTraceStore()
	loader := TraceLoader(store)

	rounds := 100
	outs := make(chan *worker.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeTraceLoaderContext(i, store)
		go func(ctx *worker.Context) {
			defer wg.Done()
			loader(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)

	assert.Len(t, outs, rounds, "TraceLoader: expected correct out count")

	for out := range outs {
		assert.Len(t, out.T.Errors, out.Keys["errors"].(int), "Marker: expected correct errors count")
		assert.Equal(t, out.Keys["uuid"].(string), out.T.GetMetadata().GetId(), "Trace metadata should be set")
	}
}
