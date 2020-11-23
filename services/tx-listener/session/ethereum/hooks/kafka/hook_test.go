// +build unit

package kafka

import (
	"context"
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"testing"

	"github.com/Shopify/sarama/mocks"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/mock"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	crc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/client/mock"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/proto"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/dynamic"
)

var c = &dynamic.Chain{
	UUID:     "test-c",
	URL:      "test-url",
	ChainID:  "888",
	Listener: &dynamic.Listener{ExternalTxEnabled: false},
}

var blockEnc = ethcommon.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0")

var block ethtypes.Block
var _ = rlp.DecodeBytes(blockEnc, &block)

func Test_AfterNewBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	receipt := &types.Receipt{
		TxHash:          "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
		ContractAddress: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
		Logs: []*types.Log{
			{
				TxHash:  "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
				Address: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
				Topics: []string{
					"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
					"0x000000000000000000000000ba826fec90cefdf6706858e5fbafcb27a290fbe0",
					"0x0000000000000000000000004aee792a88edda29932254099b9d1e06d537883f",
				},
				Data:        "0x000000000000000000000000000000000000000000000001a055690d9db80000",
				BlockNumber: block.Number().Uint64(),
				BlockHash:   block.Hash().Hex(),
				TxIndex:     0,
				Index:       0,
				Removed:     false,
			},
		},
	}
	var fakeJobs = []*entities.Job{testutils.FakeJob()}
	fakeJobs[0].Receipt = receipt

	registry := crc.NewMockContractRegistryClient(ctrl)
	ec := mock.NewMockChainStateReader(ctrl)
	txScheduler := mock2.NewMockTransactionSchedulerClient(ctrl)
	producer := mocks.NewSyncProducer(t, nil)

	// Initialize hook
	conf := &Config{
		OutTopic: "test-topic-decoded",
	}

	t.Run("should process after new block successfully", func(t *testing.T) {
		ec.EXPECT().CodeAt(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ethcommon.Hex2Bytes("0xabcd"), nil)
		registry.EXPECT().SetAccountCodeHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.SetAccountCodeHashResponse{}, nil)
		registry.EXPECT().GetEventsBySigHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.GetEventsBySigHashResponse{
			Event: "{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}",
		}, nil)
		txScheduler.EXPECT().UpdateJob(gomock.Any(), gomock.Any(), &txschedulertypes.UpdateJobRequest{
			Status:  utils.StatusMined,
			Message: fmt.Sprintf("Transaction mined in block %v", block.NumberU64()),
		}).Return(&txschedulertypes.JobResponse{}, nil)
		producer.ExpectSendMessageAndSucceed()

		expectedDecodedData := map[string]string{
			"tokens": "30000000000000000000",
			"from":   "0xBA826fEc90CEFdf6706858E5FbaFcb27A290Fbe0",
			"to":     "0x4aEE792A88eDDA29932254099b9d1e06D537883f",
		}

		hk := NewHook(conf, registry, ec, producer, txScheduler)

		var block ethtypes.Block
		err := rlp.DecodeBytes(blockEnc, &block)
		assert.NoError(t, err)

		err = hk.AfterNewBlock(context.Background(), c, &block, fakeJobs)

		assert.NoError(t, err, "AfterNewBlockEnvelope should not error")
		assert.Equal(t, "Transfer(address,address,uint256)", fakeJobs[0].Receipt.Logs[0].Event)
		assert.Equal(t, expectedDecodedData, fakeJobs[0].Receipt.Logs[0].DecodedData)
	})

	t.Run("should not fail if CodeAt' fails", func(t *testing.T) {
		ec.EXPECT().CodeAt(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("CodeAt error"))
		registry.EXPECT().GetEventsBySigHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.GetEventsBySigHashResponse{
			Event: "{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}",
		}, nil)
		txScheduler.EXPECT().UpdateJob(gomock.Any(), gomock.Any(), &txschedulertypes.UpdateJobRequest{
			Status:  utils.StatusMined,
			Message: fmt.Sprintf("Transaction mined in block %v", block.NumberU64()),
		}).Return(&txschedulertypes.JobResponse{}, nil)
		producer.ExpectSendMessageAndSucceed()

		hk := NewHook(conf, registry, ec, producer, txScheduler)
		err := hk.AfterNewBlock(context.Background(), c, &block, fakeJobs)

		assert.NoError(t, err)
	})

	t.Run("should not fail if SetAccountCodeHash fails", func(t *testing.T) {
		ec.EXPECT().CodeAt(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ethcommon.Hex2Bytes("0xabcd"), nil)
		registry.EXPECT().SetAccountCodeHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("SetAccountCodeHash error"))
		registry.EXPECT().GetEventsBySigHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.GetEventsBySigHashResponse{
			Event: "{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}",
		}, nil)
		txScheduler.EXPECT().UpdateJob(gomock.Any(), gomock.Any(), &txschedulertypes.UpdateJobRequest{
			Status:  utils.StatusMined,
			Message: fmt.Sprintf("Transaction mined in block %v", block.NumberU64()),
		}).Return(&txschedulertypes.JobResponse{}, nil)
		producer.ExpectSendMessageAndSucceed()

		hk := NewHook(conf, registry, ec, producer, txScheduler)
		err := hk.AfterNewBlock(context.Background(), c, &block, fakeJobs)

		assert.NoError(t, err)
	})

	t.Run("should not fail if GetEventsBySigHash fails", func(t *testing.T) {
		ec.EXPECT().CodeAt(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ethcommon.Hex2Bytes("0xabcd"), nil)
		registry.EXPECT().SetAccountCodeHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.SetAccountCodeHashResponse{}, nil)
		registry.EXPECT().GetEventsBySigHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error GetEventsBySigHash"))
		txScheduler.EXPECT().UpdateJob(gomock.Any(), gomock.Any(), &txschedulertypes.UpdateJobRequest{
			Status:  utils.StatusMined,
			Message: fmt.Sprintf("Transaction mined in block %v", block.NumberU64()),
		}).Return(&txschedulertypes.JobResponse{}, nil)
		producer.ExpectSendMessageAndSucceed()

		hk := NewHook(conf, registry, ec, producer, txScheduler)
		err := hk.AfterNewBlock(context.Background(), c, &block, fakeJobs)

		assert.NoError(t, err)
	})

	t.Run("should not fail if tx scheduler client fails", func(t *testing.T) {
		ec.EXPECT().CodeAt(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ethcommon.Hex2Bytes("0xabcd"), nil)
		registry.EXPECT().SetAccountCodeHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.SetAccountCodeHashResponse{}, nil)
		registry.EXPECT().GetEventsBySigHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.GetEventsBySigHashResponse{
			Event: "{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}",
		}, nil)
		txScheduler.EXPECT().UpdateJob(gomock.Any(), gomock.Any(), &txschedulertypes.UpdateJobRequest{
			Status:  utils.StatusMined,
			Message: fmt.Sprintf("Transaction mined in block %v", block.NumberU64()),
		}).Return(nil, fmt.Errorf("error"))
		producer.ExpectSendMessageAndSucceed()

		hk := NewHook(conf, registry, ec, producer, txScheduler)
		err := hk.AfterNewBlock(context.Background(), c, &block, fakeJobs)

		assert.NoError(t, err)
	})

	t.Run("should fail if it could not to produce message into kafka", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		ec.EXPECT().CodeAt(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ethcommon.Hex2Bytes("0xabcd"), nil)
		registry.EXPECT().SetAccountCodeHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.SetAccountCodeHashResponse{}, nil)
		registry.EXPECT().GetEventsBySigHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.GetEventsBySigHashResponse{
			Event: "{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}",
		}, nil)
		txScheduler.EXPECT().UpdateJob(gomock.Any(), gomock.Any(), &txschedulertypes.UpdateJobRequest{
			Status:  utils.StatusMined,
			Message: fmt.Sprintf("Transaction mined in block %v", block.NumberU64()),
		}).Return(&txschedulertypes.JobResponse{}, nil)
		producer.ExpectSendMessageAndFail(expectedErr)

		hk := NewHook(conf, registry, ec, producer, txScheduler)
		err := hk.AfterNewBlock(context.Background(), c, &block, fakeJobs)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func Test_DecodeReceipt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Initialize hook
	conf := &Config{
		OutTopic: "test-topic-decoded",
	}

	registry := crc.NewMockContractRegistryClient(ctrl)
	ec := mock.NewMockChainStateReader(ctrl)
	producer := mocks.NewSyncProducer(t, nil)
	txScheduler := mock2.NewMockTransactionSchedulerClient(ctrl)

	t.Run("should decode receipt successfully", func(t *testing.T) {
		registry.EXPECT().GetEventsBySigHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.GetEventsBySigHashResponse{
			Event: "{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}",
		}, nil)

		hk := NewHook(conf, registry, ec, producer, txScheduler)

		r := &types.Receipt{
			TxHash:          "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
			ContractAddress: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
			Logs: []*types.Log{
				{
					TxHash:  "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
					Address: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
					Topics: []string{
						"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
						"0x000000000000000000000000ba826fec90cefdf6706858e5fbafcb27a290fbe0",
						"0x0000000000000000000000004aee792a88edda29932254099b9d1e06d537883f",
					},
					Data:        "0x000000000000000000000000000000000000000000000001a055690d9db80000",
					BlockNumber: block.Number().Uint64(),
					BlockHash:   block.Hash().Hex(),
					TxIndex:     0,
					Index:       0,
					Removed:     false,
				},
			},
		}
		c := &dynamic.Chain{ChainID: "888"}
		err := hk.decodeReceipt(context.Background(), c, r)
		assert.NoError(t, err)
		assert.Equal(t, "Transfer(address,address,uint256)", r.Logs[0].Event)
		expectedDecodedData := map[string]string{
			"tokens": "30000000000000000000",
			"from":   "0xBA826fEc90CEFdf6706858E5FbaFcb27A290Fbe0",
			"to":     "0x4aEE792A88eDDA29932254099b9d1e06D537883f",
		}
		assert.Equal(t, expectedDecodedData, r.Logs[0].DecodedData)
	})

	t.Run("should not get error when not able to unmarshall event", func(t *testing.T) {
		registry.EXPECT().GetEventsBySigHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.GetEventsBySigHashResponse{
			Event: "not json event",
		}, nil)

		hk := NewHook(conf, registry, ec, producer, txScheduler)

		r := &types.Receipt{
			TxHash:          "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
			ContractAddress: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
			Logs: []*types.Log{
				{
					TxHash:  "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
					Address: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
					Topics: []string{
						"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
						"0x000000000000000000000000ba826fec90cefdf6706858e5fbafcb27a290fbe0",
						"0x0000000000000000000000004aee792a88edda29932254099b9d1e06d537883f",
					},
					Data:        "0x000000000000000000000000000000000000000000000001a055690d9db80000",
					BlockNumber: block.Number().Uint64(),
					BlockHash:   block.Hash().Hex(),
					TxIndex:     0,
					Index:       0,
					Removed:     false,
				},
			},
		}
		c := &dynamic.Chain{ChainID: "888"}
		err := hk.decodeReceipt(context.Background(), c, r)
		assert.NoError(t, err)
	})

	t.Run("should decode receipt successfully with DefaultEvents", func(t *testing.T) {
		registry.EXPECT().GetEventsBySigHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.GetEventsBySigHashResponse{
			DefaultEvents: []string{
				"{\"anonymous\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}",
				"{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}",
			},
		}, nil)

		hk := NewHook(conf, registry, ec, producer, txScheduler)

		r := &types.Receipt{
			TxHash:          "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
			ContractAddress: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
			Logs: []*types.Log{
				{
					TxHash:  "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
					Address: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
					Topics: []string{
						"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
						"0x000000000000000000000000ba826fec90cefdf6706858e5fbafcb27a290fbe0",
						"0x0000000000000000000000004aee792a88edda29932254099b9d1e06d537883f",
					},
					Data:        "0x000000000000000000000000000000000000000000000001a055690d9db80000",
					BlockNumber: block.Number().Uint64(),
					BlockHash:   block.Hash().Hex(),
					TxIndex:     0,
					Index:       0,
					Removed:     false,
				},
			},
		}
		c := &dynamic.Chain{ChainID: "888"}
		err := hk.decodeReceipt(context.Background(), c, r)
		assert.NoError(t, err)
		assert.Equal(t, "Transfer(address,address,uint256)", r.Logs[0].Event)
		expectedDecodedData := map[string]string{
			"tokens": "30000000000000000000",
			"from":   "0xBA826fEc90CEFdf6706858E5FbaFcb27A290Fbe0",
			"to":     "0x4aEE792A88eDDA29932254099b9d1e06D537883f",
		}
		assert.Equal(t, expectedDecodedData, r.Logs[0].DecodedData)
	})

	t.Run("should not fail if not finding event in DefaultEvents", func(t *testing.T) {
		registry.EXPECT().GetEventsBySigHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.GetEventsBySigHashResponse{
			DefaultEvents: []string{
				"{\"anonymous\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}",
				"{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}",
			},
		}, nil)

		hk := NewHook(conf, registry, ec, producer, txScheduler)

		r := &types.Receipt{
			TxHash:          "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
			ContractAddress: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
			Logs: []*types.Log{
				{
					TxHash:  "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
					Address: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
					Topics: []string{
						"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
						"0x000000000000000000000000ba826fec90cefdf6706858e5fbafcb27a290fbe0",
						"0x0000000000000000000000004aee792a88edda29932254099b9d1e06d537883f",
					},
					Data:        "0x000000000000000000000000000000000000000000000001a055690d9db80000",
					BlockNumber: block.Number().Uint64(),
					BlockHash:   block.Hash().Hex(),
					TxIndex:     0,
					Index:       0,
					Removed:     false,
				},
			},
		}
		c := &dynamic.Chain{ChainID: "888"}
		err := hk.decodeReceipt(context.Background(), c, r)
		assert.NoError(t, err)
	})

	t.Run("should not fail if could not unmarshal event in DefaultEvents", func(t *testing.T) {
		registry.EXPECT().GetEventsBySigHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(&contractregistry.GetEventsBySigHashResponse{
			DefaultEvents: []string{
				"could not unmarshal this event",
				"{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}",
			},
		}, nil)

		hk := NewHook(conf, registry, ec, producer, txScheduler)

		r := &types.Receipt{
			TxHash:          "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
			ContractAddress: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
			Logs: []*types.Log{
				{
					TxHash:  "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
					Address: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
					Topics: []string{
						"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
						"0x000000000000000000000000ba826fec90cefdf6706858e5fbafcb27a290fbe0",
						"0x0000000000000000000000004aee792a88edda29932254099b9d1e06d537883f",
					},
					Data:        "0x000000000000000000000000000000000000000000000001a055690d9db80000",
					BlockNumber: block.Number().Uint64(),
					BlockHash:   block.Hash().Hex(),
					TxIndex:     0,
					Index:       0,
					Removed:     false,
				},
			},
		}
		c := &dynamic.Chain{ChainID: "888"}
		err := hk.decodeReceipt(context.Background(), c, r)
		assert.NoError(t, err)
	})

	t.Run("should get an error when there are no topics", func(t *testing.T) {
		hk := NewHook(conf, registry, ec, producer, txScheduler)

		r := &types.Receipt{
			TxHash:          "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
			ContractAddress: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
			Logs: []*types.Log{
				{
					TxHash:      "0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
					Address:     "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
					Data:        "0x000000000000000000000000000000000000000000000001a055690d9db80000",
					BlockNumber: block.Number().Uint64(),
					BlockHash:   block.Hash().Hex(),
					TxIndex:     0,
					Index:       0,
					Removed:     false,
				},
			},
		}
		c := &dynamic.Chain{ChainID: "888"}
		err := hk.decodeReceipt(context.Background(), c, r)
		assert.Error(t, err)
	})
}
