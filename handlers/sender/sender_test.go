package sender

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
)

type MockTxSender struct {
	t *testing.T
}

func (s *MockTxSender) SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error {
	if chainID.Text(10) == "0" {
		return fmt.Errorf("could not send")
	}
	return nil
}

func (s *MockTxSender) SendTransaction(ctx context.Context, chainID *big.Int, args *types.SendTxArgs) (ethcommon.Hash, error) {
	if chainID.Text(10) == "0" {
		return ethcommon.Hash{}, fmt.Errorf("could not send")
	}
	return ethcommon.HexToHash("0x" + RandString(32)), nil
}

func (s *MockTxSender) SendRawPrivateTransaction(ctx context.Context, chainID *big.Int, raw string, args *types.PrivateArgs) (ethcommon.Hash, error) {
	if chainID.Text(10) == "0" {
		return ethcommon.Hash{}, fmt.Errorf("could not send")
	}
	return ethcommon.Hash{}, nil
}

var letterRunes = []rune("abcdef0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func makeSenderContext(i int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	switch i % 4 {
	case 0:
		txctx.Envelope.Chain = chain.CreateChainInt(8)
		txctx.Envelope.Tx = &ethereum.Transaction{
			Raw: ethereum.HexToData("0xabde4f3a"),
		}
		txctx.Envelope.Metadata = (&envelope.Metadata{Id: RandString(10)})
		txctx.Envelope.Tx.Hash = ethereum.HexToHash("0x" + RandString(64))
		txctx.Set("errors", 0)
		txctx.Set("status", "pending")
	case 1:
		txctx.Envelope.Chain = chain.CreateChainInt(0)
		txctx.Envelope.Tx = &ethereum.Transaction{
			Raw: ethereum.HexToData("0xabde4f3a"),
		}
		txctx.Envelope.Tx.Hash = ethereum.HexToHash("0x" + RandString(64))
		txctx.Envelope.Metadata = (&envelope.Metadata{Id: RandString(10)})
		txctx.Set("errors", 1)
		txctx.Set("status", "error")
	case 2:
		txctx.Envelope.Chain = chain.CreateChainInt(10)
		txctx.Envelope.Tx = &ethereum.Transaction{}
		txctx.Envelope.Metadata = (&envelope.Metadata{Id: RandString(10)})
		txctx.Set("errors", 0)
		txctx.Set("status", "pending")
	case 3:
		// Cannot send a transaction
		txctx.Envelope.Chain = chain.CreateChainInt(0)
		txctx.Envelope.Tx = &ethereum.Transaction{}
		txctx.Envelope.Metadata = (&envelope.Metadata{Id: RandString(10)})
		txctx.Set("errors", 1)
		txctx.Set("status", "")
	case 4:
		// Cannot send a transaction
		txctx.Envelope.Chain = chain.CreateChainInt(0)
		txctx.Envelope.Tx = &ethereum.Transaction{}
		txctx.Envelope.Metadata = (&envelope.Metadata{Id: RandString(10)})
		txctx.Set("errors", 1)
		txctx.Set("status", "")
	}
	return txctx
}

func TestSender(t *testing.T) {
	s := MockTxSender{t: t}
	store := mock.NewEnvelopeStore()
	sender := Sender(&s, store)

	rounds := 100
	outs := make(chan *engine.TxContext, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		txctx := makeSenderContext(i)
		t.Log(txctx)
		go func(txctx *engine.TxContext) {
			defer wg.Done()
			sender(txctx)
			outs <- txctx
		}(txctx)
	}
	wg.Wait()
	close(outs)

	assert.Len(t, outs, rounds, "Marker: expected correct out count")

	for out := range outs {
		assert.Len(t, out.Envelope.Errors, out.Get("errors").(int), "Marker: expected correct errors count")
		status, _, _ := store.GetStatus(context.Background(), out.Envelope.GetMetadata().GetId())
		assert.Equal(t, out.Get("status").(string), status, "Transaction should be in status %q", "pending")
	}
}
