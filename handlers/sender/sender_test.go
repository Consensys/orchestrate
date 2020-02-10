package sender

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"

	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store/client/mocks"
)

type MockTxSender struct {
	t *testing.T
}

const (
	endpointNoError = "testURL"
	endpointError   = "error"
)

func (s *MockTxSender) SendRawTransaction(ctx context.Context, endpoint, raw string) error {
	if endpoint == endpointError {
		return fmt.Errorf("mock: failed to send a raw transaction")
	}
	return nil
}

func (s *MockTxSender) SendTransaction(ctx context.Context, endpoint string, args *types.SendTxArgs) (ethcommon.Hash, error) {
	if endpoint == endpointError {
		return ethcommon.Hash{}, fmt.Errorf("mock: failed to send an unsigned transaction")
	}
	return ethcommon.HexToHash("0x" + RandString(32)), nil
}

func (s *MockTxSender) SendRawPrivateTransaction(ctx context.Context, endpoint, raw string, args *types.PrivateArgs) (ethcommon.Hash, error) {
	if endpoint == endpointError {
		return ethcommon.Hash{}, fmt.Errorf("mock: failed to send a raw private transaction")
	}
	return ethcommon.Hash{}, nil
}

func (s *MockTxSender) SendQuorumRawPrivateTransaction(ctx context.Context, endpoint, signedTxHash string, privateFor []string) (ethcommon.Hash, error) {
	if endpoint == endpointError {
		return ethcommon.Hash{}, fmt.Errorf("mock: failed to send a raw Tessera transaction")
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
	txRaw := "0xabde4f3a"
	txHash := "0x" + RandString(64)
	switch i % 8 {
	case 0:
		// Valid send base transaction
		txctx.WithContext(proxy.With(txctx.Context(), endpointNoError))
		_ = txctx.Builder.
			SetID(RandString(10)).
			SetTxHash(ethcommon.HexToHash(txHash)).
			SetRawString(txRaw)
		txctx.Set("status", "PENDING")
	case 1:
		// Invalid send base transaction
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Builder.
			SetID(RandString(10)).
			SetTxHash(ethcommon.HexToHash(txHash)).
			SetRawString(txRaw)
		txctx.Set("error", "mock: failed to send a raw transaction")
		txctx.Set("status", "ERROR")
	case 2:
		//
		txctx.WithContext(proxy.With(txctx.Context(), endpointNoError))
		_ = txctx.Builder.SetID(RandString(10)).SetRawString(txRaw)
		txctx.Set("status", "PENDING")
	case 3:
		// Cannot send a public transaction
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Builder.SetID(RandString(10))
		txctx.Set("error", "no raw filled")
		txctx.Set("status", "")
	case 4:
		// Cannot send a Besu Orion transaction
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Builder.SetID(RandString(10)).SetMethod(tx.Method_EEA_SENDPRIVATETRANSACTION).SetRawString(txRaw)
		txctx.Set("error", "mock: failed to send a raw private transaction")
		txctx.Set("status", "")
	case 5:
		// Cannot send a Quorum Tessera transaction
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Builder.SetID(RandString(10)).SetMethod(tx.Method_ETH_SENDRAWPRIVATETRANSACTION).MustSetFromString("0x1").SetPrivateFor([]string{"test"}).SetRawString(txRaw)
		txctx.Set("error", "mock: failed to send a raw Tessera transaction")
		txctx.Set("status", "ERROR")
	case 6:
		// Cannot send a Quorum Constellation transaction
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Builder.SetID(RandString(10)).SetMethod(tx.Method_ETH_SENDPRIVATETRANSACTION).MustSetFromString("0x1").SetRawString(txRaw)
		txctx.Set("error", "mock: failed to send an unsigned transaction")
		txctx.Set("status", "")
	case 7:
		// 	// Cannot send a transaction with unknown protocol type
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Builder.SetID(RandString(10)).SetMethod(123).MustSetFromString("0x1").SetRawString(txRaw)
		txctx.Set("error", "invalid private protocol \"123\"")
		txctx.Set("status", "")
	case 8:
		// Cannot send a signed private transaction with Constellation protocol
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Builder.
			SetID(RandString(10)).
			SetMethod(tx.Method_ETH_SENDPRIVATETRANSACTION).
			SetTxHash(ethcommon.HexToHash(txHash)).
			SetRawString(txRaw)
		txctx.Set("error", "mock: failed to send an unsigned transaction")
		txctx.Set("status", "")
	}
	return txctx
}

func TestSender(t *testing.T) {
	s := MockTxSender{t: t}
	client := clientmock.New()
	sender := Sender(&s, client)

	rounds := 15
	outs := make(chan *engine.TxContext, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		txctx := makeSenderContext(i)

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
		resp, _ := client.LoadByID(
			context.Background(),
			&evlpstore.LoadByIDRequest{
				Id: out.Builder.GetID(),
			},
		)

		expectedError := out.Get("error")
		if expectedError != nil {
			assert.Equal(t, expectedError.(string), out.Builder.Errors[0].Message, "")
		} else {
			assert.Equal(t, out.Get("status").(string), resp.GetStatusInfo().GetStatus().String(), "Incorrect envelope status")
			assert.Len(t, out.Builder.Errors, 0, "")
		}
	}
}
