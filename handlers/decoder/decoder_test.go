package decoder

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry/client/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

const (
	testEvent = "Transfer(address,address,uint256)"
)

func TestDecoder(t *testing.T) {
	testSet := []struct {
		name          string
		input         func(txctx *engine.TxContext) *engine.TxContext
		expectedTxctx func(txctx *engine.TxContext) *engine.TxContext
	}{
		{
			"Receipt without error and log decoded",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Builder.SetChainID(big.NewInt(1))
				txctx.Builder.Receipt = &ethereum.Receipt{
					Logs: []*ethereum.Log{
						{
							Data: "0x000000000000000000000000000000000000000000000001a055690d9db80000",
							Topics: []string{
								"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
								"0x000000000000000000000000ba826fec90cefdf6706858e5fbafcb27a290fbe0",
								"0x0000000000000000000000004aee792a88edda29932254099b9d1e06d537883f",
							},
						},
					},
				}
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Builder.GetReceipt().Logs[0].DecodedData = map[string]string{
					"from":   "0xBA826fEc90CEFdf6706858E5FbaFcb27A290Fbe0",
					"to":     "0x4aEE792A88eDDA29932254099b9d1e06D537883f",
					"tokens": "30000000000000000000",
				}
				txctx.Builder.GetReceipt().Logs[0].Event = testEvent
				return txctx
			},
		},
		{
			"Receipt without error and unknown abi",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Builder.SetChainID(big.NewInt(1))
				txctx.Builder.Receipt = &ethereum.Receipt{
					Logs: []*ethereum.Log{
						{
							Data: "0x000000000000000000000000000000000000000000000001a055690d9db80000",
							Topics: []string{
								"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
								"0x000000000000000000000000ba826fec90cefdf6706858e5fbafcb27a290fbe0",
								"0x0000000000000000000000004aee792a88edda29932254099b9d1e06d537883f",
							},
						},
					},
				}
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Builder.GetReceipt().Logs[0].DecodedData = map[string]string{
					"from":   "0xBA826fEc90CEFdf6706858E5FbaFcb27A290Fbe0",
					"to":     "0x4aEE792A88eDDA29932254099b9d1e06D537883f",
					"tokens": "30000000000000000000",
				}
				txctx.Builder.GetReceipt().Logs[0].Event = testEvent
				return txctx
			},
		},
		{
			"Receipt without topics",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Builder.SetChainID(big.NewInt(1))
				txctx.Builder.Receipt = &ethereum.Receipt{
					Logs: []*ethereum.Log{
						{
							Data:   "0x000000000000000000000000000000000000000000000001a055690d9db80000",
							Topics: []string{},
						},
					},
				}
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.InternalError("invalid receipt (no topics in log)").ExtendComponent(component)
				txctx.Builder.Errors = append(txctx.Builder.Errors, err)
				return txctx
			},
		},
	}

	var testABI = `[
		{
			"anonymous":false,
			"inputs":[
				{"indexed":true,"name":"from","type":"address"},
				{"indexed":true,"name":"to","type":"address"},
				{"indexed":false,"name":"tokens","type":"uint256"}
			],
			"name":"Transfer",
			"type":"event"
		}
	]`

	registry := clientmock.New()
	_, err := registry.RegisterContract(context.Background(),
		&contractregistry.RegisterContractRequest{
			Contract: &abi.Contract{
				Id: &abi.ContractId{
					Name: "known",
				},
				Abi:              testABI,
				Bytecode:         hexutil.Encode([]byte{1, 2, 3}),
				DeployedBytecode: hexutil.Encode([]byte{1, 2}),
			},
		})
	if err != nil {
		t.Fatalf("could not register contract - %v", err)
	}
	h := Decoder(registry)

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			h(test.input(txctx))

			expectedTxctx := engine.NewTxContext()
			expectedTxctx.Logger = txctx.Logger
			expectedTxctx = test.expectedTxctx(test.input(expectedTxctx))

			assert.True(t, reflect.DeepEqual(txctx, expectedTxctx), "Expected same input")
		})
	}
}
