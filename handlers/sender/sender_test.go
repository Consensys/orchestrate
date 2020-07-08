// +build unit

package sender

import (
	"context"
	"fmt"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client/mock"
	"math/rand"
	"sync"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/client/mock"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

const (
	endpointNoError = "testURL"
	endpointError   = "error"
)

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
		_ = txctx.Envelope.
			SetID(RandString(10)).
			SetTxHash(ethcommon.HexToHash(txHash)).
			SetRawString(txRaw)
		txctx.Set("status", "PENDING")
	case 1:
		// Invalid send base transaction
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Envelope.
			SetID(RandString(10)).
			SetTxHash(ethcommon.HexToHash(txHash)).
			SetRawString(txRaw)
		txctx.Set("error", "mock: failed to send a raw transaction")
		txctx.Set("status", "ERROR")
	case 2:
		//
		txctx.WithContext(proxy.With(txctx.Context(), endpointNoError))
		_ = txctx.Envelope.SetID(RandString(10)).SetRawString(txRaw)
		txctx.Set("status", "PENDING")
	case 3:
		// Cannot send a public transaction
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Envelope.SetID(RandString(10))
		txctx.Set("error", "no raw filled")
		txctx.Set("status", "")
	case 4:
		// Cannot send a Besu Orion transaction
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Envelope.SetID(RandString(10)).SetMethod(tx.Method_EEA_SENDPRIVATETRANSACTION).SetRawString(txRaw)
		txctx.Set("error", "mock: failed to send a raw private transaction")
		txctx.Set("status", "")
	case 5:
		// Cannot send a Quorum Tessera transaction
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Envelope.SetID(RandString(10)).SetMethod(tx.Method_ETH_SENDRAWPRIVATETRANSACTION).MustSetFromString("0x1").SetPrivateFor([]string{"test"}).SetRawString(txRaw)
		txctx.Set("error", "mock: failed to send a raw Tessera transaction")
		txctx.Set("status", "ERROR")
	case 6:
		// Cannot send a Quorum Constellation transaction
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Envelope.SetID(RandString(10)).SetMethod(tx.Method_ETH_SENDPRIVATETRANSACTION).MustSetFromString("0x1").SetRawString(txRaw)
		txctx.Set("error", "mock: failed to send an unsigned transaction")
		txctx.Set("status", "")
	case 7:
		// 	// Cannot send a transaction with unknown protocol type
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Envelope.SetID(RandString(10)).SetMethod(123).MustSetFromString("0x1").SetRawString(txRaw)
		txctx.Set("error", "invalid transaction protocol \"123\"")
		txctx.Set("status", "")
	case 8:
		// Cannot send a signed private transaction with Constellation protocol
		txctx.WithContext(proxy.With(txctx.Context(), endpointError))
		_ = txctx.Envelope.
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
	ctrl := gomock.NewController(t)
	client := clientmock.NewMockEnvelopeStoreClient(ctrl)
	txSchedulerClient := mock2.NewMockTransactionSchedulerClient(ctrl)
	client.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&svc.StoreRequest{}), gomock.Any()).Times(15)
	client.EXPECT().SetStatus(gomock.Any(), gomock.AssignableToTypeOf(&svc.SetStatusRequest{})).Times(15)
	client.EXPECT().LoadByID(gomock.Any(), gomock.AssignableToTypeOf(&svc.LoadByIDRequest{})).
		Times(15).Return(&svc.StoreResponse{
		StatusInfo: &svc.StatusInfo{Status: svc.Status_PENDING},
	}, nil)

	s := mock.NewMockTransactionSender(ctrl)
	s.EXPECT().SendQuorumRawPrivateTransaction(gomock.Any(), gomock.Eq(endpointError), gomock.Any(), gomock.Any()).Return(ethcommon.Hash{}, fmt.Errorf("mock: failed to send a raw Tessera transaction")).AnyTimes()
	s.EXPECT().SendQuorumRawPrivateTransaction(gomock.Any(), gomock.Not(gomock.Eq(endpointError)), gomock.Any(), gomock.Any()).Return(ethcommon.Hash{}, nil).AnyTimes()
	s.EXPECT().SendRawPrivateTransaction(gomock.Any(), gomock.Eq(endpointError), gomock.Any()).Return(ethcommon.Hash{}, fmt.Errorf("mock: failed to send a raw private transaction")).AnyTimes()
	s.EXPECT().SendRawPrivateTransaction(gomock.Any(), gomock.Not(gomock.Eq(endpointError)), gomock.Any()).Return(ethcommon.Hash{}, nil).AnyTimes()
	s.EXPECT().SendTransaction(gomock.Any(), gomock.Eq(endpointError), gomock.Any()).Return(ethcommon.Hash{}, fmt.Errorf("mock: failed to send an unsigned transaction")).AnyTimes()
	s.EXPECT().SendTransaction(gomock.Any(), gomock.Not(gomock.Eq(endpointError)), gomock.Any()).Return(ethcommon.HexToHash("0x"+RandString(32)), nil).AnyTimes()
	s.EXPECT().SendRawTransaction(gomock.Any(), gomock.Eq(endpointError), gomock.Any()).Return(fmt.Errorf("mock: failed to send a raw transaction")).AnyTimes()
	s.EXPECT().SendRawTransaction(gomock.Any(), gomock.Not(gomock.Eq(endpointError)), gomock.Any()).Return(nil).AnyTimes()
	sender := Sender(s, client, txSchedulerClient)

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
			&svc.LoadByIDRequest{
				Id: out.Envelope.GetID(),
			},
		)

		expectedError := out.Get("error")
		if expectedError != nil {
			assert.Equal(t, expectedError.(string), out.Envelope.Errors[0].Message, "")
		} else {
			assert.Equal(t, out.Get("status").(string), resp.GetStatusInfo().GetStatus().String(), "Incorrect envelope status")
			assert.Len(t, out.Envelope.Errors, 0, "")
		}
	}
}
