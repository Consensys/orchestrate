package handlers

import (
	"context"
	"fmt"
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

type MockTxSender struct {
	t *testing.T
}

func (s *MockTxSender) SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error {
	if chainID.Text(10) == "0" {
		return fmt.Errorf("Could not send")
	}
	return nil
}

var letterRunes = []rune("abcdefgABCDEF0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func makeSenderContext(i int) *worker.Context {
	ctx := worker.NewContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	switch i % 4 {
	case 0:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(10))
		ctx.T.Tx = (&ethereum.Transaction{}).SetRaw("0xabde4f3a")
		ctx.T.Metadata = (&trace.Metadata{Id: RandString(10)})
		ctx.T.Tx.Hash = "0x" + RandString(32)
		ctx.Keys["errors"] = 0
		ctx.Keys["status"] = "pending"
	case 1:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(0))
		ctx.T.Tx = (&ethereum.Transaction{}).SetRaw("0xabde4f3a")
		ctx.T.Tx.Hash = "0x" + RandString(32)
		ctx.T.Metadata = (&trace.Metadata{Id: RandString(10)})
		ctx.Keys["errors"] = 1
		ctx.Keys["status"] = "error"
	case 2:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(0))
		ctx.T.Tx = (&ethereum.Transaction{}).SetRaw(``)
		ctx.T.Tx.Hash = "0x" + RandString(32)
		ctx.T.Metadata = (&trace.Metadata{Id: RandString(10)})
		ctx.Keys["errors"] = 0
		ctx.Keys["status"] = ""
	case 3:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(10))
		ctx.T.Tx = (&ethereum.Transaction{}).SetRaw(``)
		ctx.T.Tx.Hash = "0x" + RandString(32)
		ctx.T.Metadata = (&trace.Metadata{Id: RandString(10)})
		ctx.Keys["errors"] = 0
		ctx.Keys["status"] = ""
	}
	return ctx
}

func TestSender(t *testing.T) {
	s := MockTxSender{t: t}
	store := mock.NewTraceStore()
	sender := Sender(&s, store)

	rounds := 100
	outs := make(chan *worker.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeSenderContext(i)
		go func(ctx *worker.Context) {
			defer wg.Done()
			sender(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)

	assert.Len(t, outs, rounds, "Marker: expected correct out count")

	for out := range outs {
		assert.Len(t, out.T.Errors, out.Keys["errors"].(int), "Marker: expected correct errors count")
		status, _, _ := store.GetStatus(context.Background(), out.T.GetMetadata().GetId())
		assert.Equal(t, out.Keys["status"].(string), status, "Transaction should be in status %q", "pending")
	}
}
