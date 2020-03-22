// +build unit

package sarama

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Shopify/sarama/mocks"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types/testutils"
)

func TestCredit(t *testing.T) {
	testSet := []struct {
		name                     string
		ctx                      context.Context
		req                      testutils.TestRequest
		expectedToSendTxEnvelope bool
		expectedTxEnvelopeSent   *tx.Envelope
		sendMesssageError        error
	}{
		{
			"credit without error",
			context.Background(),
			testutils.TestRequest{
				Req: &types.Request{
					ParentTxID:  "testID",
					ChainID:     big.NewInt(10),
					ChainURL:    "testURL",
					ChainName:   "testChainName",
					ChainUUID:   "testChainUUID",
					Beneficiary: ethcommon.HexToAddress("0xcd"),
					FaucetsCandidates: map[string]types.Faucet{
						"test": {
							Amount:     big.NewInt(1),
							MaxBalance: big.NewInt(20),
							Cooldown:   time.Second,
							Creditor:   ethcommon.HexToAddress("0xab"),
						},
						"supposedToBeElected": {
							Amount:     big.NewInt(2),
							MaxBalance: big.NewInt(20),
							Cooldown:   time.Second,
							Creditor:   ethcommon.HexToAddress("0xab"),
						},
					},
				},
				ExpectedAmount: big.NewInt(2),
			},
			true,
			tx.NewEnvelope().
				SetContextLabelsValue("faucet.parentTxID", "testID").
				SetInternalLabelsValue("chainID", "10").
				SetInternalLabelsValue("chainUUID", "testChainUUID").
				SetFrom(ethcommon.HexToAddress("0xab")).
				SetTo(ethcommon.HexToAddress("0xcd")).
				SetValue(big.NewInt(2)).
				SetChainID(big.NewInt(10)).
				SetChainName("testChainName").
				SetChainUUID("testChainUUID"),
			nil,
		},
		{
			"credit without faucet candidates",
			context.Background(),
			testutils.TestRequest{
				Req: &types.Request{
					ParentTxID:  "testID",
					ChainID:     big.NewInt(10),
					ChainURL:    "testURL",
					ChainName:   "testChainName",
					ChainUUID:   "testChainUUID",
					Beneficiary: ethcommon.HexToAddress("0xcd"),
				},
				ExpectedErr: errors.FaucetWarning("no faucet request").ExtendComponent(component),
			},
			false,
			nil,
			nil,
		},
		{
			"credit without error",
			context.Background(),
			testutils.TestRequest{
				Req: &types.Request{
					ParentTxID:  "testID",
					ChainID:     big.NewInt(10),
					ChainURL:    "testURL",
					ChainName:   "testChainName",
					ChainUUID:   "testChainUUID",
					Beneficiary: ethcommon.HexToAddress("0xcd"),
					FaucetsCandidates: map[string]types.Faucet{
						"test": {
							Amount:     big.NewInt(1),
							MaxBalance: big.NewInt(20),
							Cooldown:   time.Second,
							Creditor:   ethcommon.HexToAddress("0xab"),
						},
						"supposedToBeElected": {
							Amount:     big.NewInt(2),
							MaxBalance: big.NewInt(20),
							Cooldown:   time.Second,
							Creditor:   ethcommon.HexToAddress("0xab"),
						},
					},
				},
				ExpectedErr: errors.KafkaConnectionError("could not send faucet transaction - got kafka send message error").ExtendComponent(component),
			},
			true,
			tx.NewEnvelope().
				SetContextLabelsValue("faucet.parentTxID", "testID").
				SetInternalLabelsValue("chainID", "10").
				SetInternalLabelsValue("chainUUID", "testChainUUID").
				SetFrom(ethcommon.HexToAddress("0xab")).
				SetTo(ethcommon.HexToAddress("0xcd")).
				SetValue(big.NewInt(2)).
				SetChainID(big.NewInt(10)).
				SetChainName("testChainName").
				SetChainUUID("testChainUUID"),
			fmt.Errorf("kafka send message error"),
		},
	}

	p := mocks.NewSyncProducer(t, nil)
	f := NewFaucet(p)

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.expectedToSendTxEnvelope {
				valueChecker := func(val []byte) error {
					txEnvelope := &tx.TxEnvelope{}
					err := encoding.Unmarshal(val, txEnvelope)
					if err != nil {
						return err
					}
					envelope, err := txEnvelope.Envelope()
					if err != nil {
						return err
					}

					assert.NotEmpty(t, envelope.GetID())
					_ = envelope.SetID("")
					assert.Equal(t, test.expectedTxEnvelopeSent, envelope)
					return nil
				}

				if test.sendMesssageError != nil {
					p.ExpectSendMessageWithCheckerFunctionAndFail(valueChecker, test.sendMesssageError)
				} else {
					p.ExpectSendMessageWithCheckerFunctionAndSucceed(valueChecker)
				}
			}

			test.req.ResultAmount, test.req.ResultErr = f.Credit(context.Background(), test.req.Req)
			testutils.AssertRequest(t, &test.req)
		})
	}

}
