//+build unit
//+build !race

package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	backoffmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/backoff/mock"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	mock3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/session/ethereum/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock6 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/mock"
	mock4 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/dynamic"
	mock5 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/metrics/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/session/ethereum/hooks/mock"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/session/ethereum/offset/mock"
)

func TestSession_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	blockPosition := big.NewInt(1)
	eeaPrivPrecompiledContractAddr := "0x000000000000000000000000000000000000007E"
	newBlockPosition := big.NewInt(blockPosition.Int64() + 1)
	receipt := newFakeReceipt()
	toAddress := "0x0000000000000000000000000000000000000001"
	txHash := "0xfda31b00fbfc77c8aaaf225ff05098324fcf2e2515d95b488022b42b3b946144"
	txHashPrivate := "0xa99cc8da7063b04f7e350da806db7e6aa92e292b562d7a5918dcf8c02a9a2aea"

	mockHook := mock.NewMockHook(ctrl)
	mockOffsetManager := mock2.NewMockManager(ctrl)
	mockEthClient := mock3.NewMockEthClient(ctrl)
	mockClient := mock4.NewMockOrchestrateClient(ctrl)
	mockMetrics := mock5.NewMockListenerMetrics(ctrl)

	blockCounter := mock6.NewMockCounter(ctrl)
	blockCounter.EXPECT().With(gomock.Any()).AnyTimes().Return(blockCounter)
	blockCounter.EXPECT().Add(gomock.Any()).AnyTimes()
	mockMetrics.EXPECT().BlockCounter().AnyTimes().Return(blockCounter)

	t.Run("should process block successfully with internal txs", func(t *testing.T) {
		cancellableCtx, cancel := context.WithCancel(ctx)
		jobResponse := testutils.FakeJobResponse()
		chain := newFakeChain()
		block := newFakeBlock(newBlockPosition, toAddress)
		jobResponse.Transaction.Hash = txHash
		session := NewSession(chain, mockEthClient, mockClient, mockHook, mockOffsetManager, mockMetrics)
		bckoff := &backoffmock.MockIntervalBackoff{}
		session.bckOff = bckoff

		mockOffsetManager.EXPECT().GetLastBlockNumber(gomock.Any(), chain).Return(blockPosition.Uint64(), nil)
		mockEthClient.EXPECT().EEAPrivPrecompiledContractAddr(gomock.Any(), chain.URL).Return(common.HexToAddress(eeaPrivPrecompiledContractAddr), nil)
		mockEthClient.EXPECT().HeaderByNumber(gomock.Any(), chain.URL, nil).Return(&types.Header{
			Number: newBlockPosition,
		}, nil).AnyTimes()
		mockEthClient.EXPECT().BlockByNumber(gomock.Any(), chain.URL, newBlockPosition).Return(block, nil)
		mockClient.EXPECT().
			SearchJob(gomock.Any(), &entities.JobFilters{
				TxHashes:  []string{txHash},
				ChainUUID: chain.UUID,
				Status:    entities.StatusPending,
			}).
			Return([]*txschedulertypes.JobResponse{jobResponse}, nil)
		mockEthClient.EXPECT().TransactionReceipt(gomock.Any(), chain.URL, common.HexToHash(txHash)).Return(receipt, nil)
		mockHook.EXPECT().AfterNewBlock(gomock.Any(), chain, block, gomock.Any()).Return(nil)
		mockOffsetManager.EXPECT().SetLastBlockNumber(gomock.Any(), chain, newBlockPosition.Uint64()).Return(nil)

		// Start session
		exitErr := make(chan error)
		go func() {
			exitErr <- session.Run(cancellableCtx)
		}()

		go func() {
			<-time.After(200 * time.Microsecond)
			cancel()
		}()

		select {
		case <-exitErr:
			assert.False(t, bckoff.HasRetried())
		// Inject hook error
		case <-time.After(1 * time.Second):
			assert.Fail(t, "should have finished")
		}
	})

	t.Run("should process block successfully with internal txs in batches", func(t *testing.T) {
		cancellableCtx, cancel := context.WithCancel(ctx)

		// Generate 70 txs that will be processed in 3 batches of 30, 30 and 10
		var txHashes []string
		var txs []*types.Transaction
		for i := 0; i < 70; i++ {
			tx := types.NewTransaction(
				0,
				common.HexToAddress(toAddress),
				big.NewInt(10000),
				uint64(21000),
				big.NewInt(20000),
				[]byte{},
			)

			txHashes = append(txHashes, txHash)
			txs = append(txs, tx)
		}

		block := types.NewBlock(&types.Header{Number: blockPosition}, txs, []*types.Header{}, []*types.Receipt{})
		chain := newFakeChain()
		session := NewSession(chain, mockEthClient, mockClient, mockHook, mockOffsetManager, mockMetrics)
		backoff := &backoffmock.MockIntervalBackoff{}
		session.bckOff = backoff

		mockOffsetManager.EXPECT().GetLastBlockNumber(gomock.Any(), chain).Return(blockPosition.Uint64(), nil)
		mockEthClient.EXPECT().EEAPrivPrecompiledContractAddr(gomock.Any(), chain.URL).Return(common.HexToAddress(eeaPrivPrecompiledContractAddr), nil)
		mockEthClient.EXPECT().HeaderByNumber(gomock.Any(), chain.URL, nil).Return(&types.Header{
			Number: newBlockPosition,
		}, nil).AnyTimes()
		mockEthClient.EXPECT().BlockByNumber(gomock.Any(), chain.URL, newBlockPosition).Return(block, nil)
		mockClient.EXPECT().
			SearchJob(gomock.Any(), &entities.JobFilters{
				TxHashes:  txHashes[0:MaxTxHashesLength],
				ChainUUID: chain.UUID,
				Status:    entities.StatusPending,
			}).
			Return([]*txschedulertypes.JobResponse{}, nil)
		mockClient.EXPECT().
			SearchJob(gomock.Any(), &entities.JobFilters{
				TxHashes:  txHashes[MaxTxHashesLength : MaxTxHashesLength*2],
				ChainUUID: chain.UUID,
				Status:    entities.StatusPending,
			}).
			Return([]*txschedulertypes.JobResponse{}, nil)
		mockClient.EXPECT().
			SearchJob(gomock.Any(), &entities.JobFilters{
				TxHashes:  txHashes[MaxTxHashesLength*2 : 70],
				ChainUUID: chain.UUID,
				Status:    entities.StatusPending,
			}).
			Return([]*txschedulertypes.JobResponse{}, nil)
		mockHook.EXPECT().AfterNewBlock(gomock.Any(), chain, block, gomock.Any()).Return(nil)
		mockOffsetManager.EXPECT().SetLastBlockNumber(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		// Start session
		exitErr := make(chan error)
		go func() {
			exitErr <- session.Run(cancellableCtx)
		}()

		go func() {
			<-time.After(200 * time.Microsecond)
			cancel()
		}()

		select {
		case <-exitErr:
			assert.False(t, backoff.HasRetried())
		// Inject hook error
		case <-time.After(1 * time.Second):
			assert.Fail(t, "should have finished")
		}
	})

	t.Run("should process block successfully with internal private txs", func(t *testing.T) {
		cancellableCtx, cancel := context.WithCancel(ctx)
		jobResponse := testutils.FakeJobResponse()
		chain := newFakeChain()
		block := newFakeBlock(newBlockPosition, eeaPrivPrecompiledContractAddr)
		jobResponse.Transaction.Hash = txHashPrivate
		session := NewSession(chain, mockEthClient, mockClient, mockHook, mockOffsetManager, mockMetrics)
		bckoff := &backoffmock.MockIntervalBackoff{}
		session.bckOff = bckoff

		mockOffsetManager.EXPECT().GetLastBlockNumber(gomock.Any(), chain).Return(blockPosition.Uint64(), nil)
		mockEthClient.EXPECT().EEAPrivPrecompiledContractAddr(gomock.Any(), chain.URL).Return(common.HexToAddress(eeaPrivPrecompiledContractAddr), nil)
		mockEthClient.EXPECT().HeaderByNumber(gomock.Any(), chain.URL, nil).Return(&types.Header{
			Number: newBlockPosition,
		}, nil).AnyTimes()
		mockEthClient.EXPECT().BlockByNumber(gomock.Any(), chain.URL, newBlockPosition).Return(block, nil)
		mockClient.EXPECT().
			SearchJob(gomock.Any(), &entities.JobFilters{
				TxHashes:  []string{txHashPrivate},
				ChainUUID: chain.UUID,
				Status:    entities.StatusPending,
			}).
			Return([]*txschedulertypes.JobResponse{jobResponse}, nil)
		mockEthClient.EXPECT().PrivateTransactionReceipt(gomock.Any(), chain.URL, common.HexToHash(txHashPrivate)).Return(receipt, nil)
		mockHook.EXPECT().AfterNewBlock(gomock.Any(), chain, block, gomock.Any()).Return(nil)
		mockOffsetManager.EXPECT().SetLastBlockNumber(gomock.Any(), chain, newBlockPosition.Uint64()).Return(nil)

		// Start session
		exitErr := make(chan error)
		go func() {
			exitErr <- session.Run(cancellableCtx)
		}()

		go func() {
			<-time.After(200 * time.Microsecond)
			cancel()
		}()

		select {
		case <-exitErr:
			assert.False(t, bckoff.HasRetried())
		// Inject hook error
		case <-time.After(1 * time.Second):
			assert.Fail(t, "should have finished")
		}
	})

	t.Run("should fetch receipts successfully for external txs", func(t *testing.T) {
		cancellableCtx, cancel := context.WithCancel(ctx)
		block := newFakeBlock(newBlockPosition, toAddress)
		chain := newFakeChain()
		chain.Listener.ExternalTxEnabled = true
		session := NewSession(chain, mockEthClient, mockClient, mockHook, mockOffsetManager, mockMetrics)
		bckoff := &backoffmock.MockIntervalBackoff{}
		session.bckOff = bckoff

		mockOffsetManager.EXPECT().GetLastBlockNumber(gomock.Any(), chain).Return(blockPosition.Uint64(), nil)
		mockEthClient.EXPECT().EEAPrivPrecompiledContractAddr(gomock.Any(), chain.URL).Return(common.HexToAddress(eeaPrivPrecompiledContractAddr), nil)
		mockEthClient.EXPECT().HeaderByNumber(gomock.Any(), chain.URL, nil).Return(&types.Header{
			Number: newBlockPosition,
		}, nil).AnyTimes()
		mockEthClient.EXPECT().BlockByNumber(gomock.Any(), chain.URL, newBlockPosition).Return(block, nil)
		mockClient.EXPECT().
			SearchJob(gomock.Any(), &entities.JobFilters{
				TxHashes:  []string{txHash},
				ChainUUID: chain.UUID,
				Status:    entities.StatusPending,
			}).
			Return([]*txschedulertypes.JobResponse{}, nil)
		mockEthClient.EXPECT().TransactionReceipt(gomock.Any(), chain.URL, common.HexToHash(txHash)).Return(receipt, nil)
		mockHook.EXPECT().AfterNewBlock(gomock.Any(), chain, block, gomock.Any()).Return(nil)
		mockOffsetManager.EXPECT().SetLastBlockNumber(gomock.Any(), chain, newBlockPosition.Uint64()).Return(nil)

		// Start session
		exitErr := make(chan error)
		go func() {
			exitErr <- session.Run(cancellableCtx)
		}()

		go func() {
			<-time.After(200 * time.Microsecond)
			cancel()
		}()

		select {
		case <-exitErr:
			assert.False(t, bckoff.HasRetried())
		// Inject hook error
		case <-time.After(1 * time.Second):
			assert.Fail(t, "should have finished")
		}
	})

	t.Run("should fetch receipts successfully for external private txs", func(t *testing.T) {
		cancellableCtx, cancel := context.WithCancel(ctx)
		block := newFakeBlock(newBlockPosition, eeaPrivPrecompiledContractAddr)
		chain := newFakeChain()
		chain.Listener.ExternalTxEnabled = true
		session := NewSession(chain, mockEthClient, mockClient, mockHook, mockOffsetManager, mockMetrics)
		bckoff := &backoffmock.MockIntervalBackoff{}
		session.bckOff = bckoff

		mockOffsetManager.EXPECT().GetLastBlockNumber(gomock.Any(), chain).Return(blockPosition.Uint64(), nil)
		mockEthClient.EXPECT().EEAPrivPrecompiledContractAddr(gomock.Any(), chain.URL).Return(common.HexToAddress(eeaPrivPrecompiledContractAddr), nil)
		mockEthClient.EXPECT().HeaderByNumber(gomock.Any(), chain.URL, nil).Return(&types.Header{
			Number: newBlockPosition,
		}, nil).AnyTimes()
		mockEthClient.EXPECT().BlockByNumber(gomock.Any(), chain.URL, newBlockPosition).Return(block, nil)
		mockClient.EXPECT().
			SearchJob(gomock.Any(), &entities.JobFilters{
				TxHashes:  []string{txHashPrivate},
				ChainUUID: chain.UUID,
				Status:    entities.StatusPending,
			}).
			Return([]*txschedulertypes.JobResponse{}, nil)
		mockEthClient.EXPECT().PrivateTransactionReceipt(gomock.Any(), chain.URL, common.HexToHash(txHashPrivate)).Return(receipt, nil)
		mockHook.EXPECT().AfterNewBlock(gomock.Any(), chain, block, gomock.Any()).Return(nil)
		mockOffsetManager.EXPECT().SetLastBlockNumber(gomock.Any(), chain, newBlockPosition.Uint64()).Return(nil)

		// Start session
		exitErr := make(chan error)
		go func() {
			exitErr <- session.Run(cancellableCtx)
		}()

		go func() {
			<-time.After(200 * time.Microsecond)
			cancel()
		}()

		select {
		case <-exitErr:
			assert.False(t, bckoff.HasRetried())
		// Inject hook error
		case <-time.After(1 * time.Second):
			assert.Fail(t, "should have finished")
		}
	})

	t.Run("should fail and retry if GetLastBlockNumber fails", func(t *testing.T) {
		cancellableCtx, cancel := context.WithCancel(ctx)
		chain := newFakeChain()
		session := NewSession(chain, mockEthClient, mockClient, mockHook, mockOffsetManager, mockMetrics)
		bckoff := &backoffmock.MockIntervalBackoff{}
		session.bckOff = bckoff

		mockOffsetManager.EXPECT().GetLastBlockNumber(gomock.Any(), chain).Return(uint64(0), fmt.Errorf("error")).AnyTimes()

		// Start session
		exitErr := make(chan error)
		go func() {
			exitErr <- session.Run(cancellableCtx)
		}()
		go func() {
			<-time.After(200 * time.Microsecond)
			cancel()
		}()

		select {
		case <-exitErr:
			assert.True(t, bckoff.HasRetried())
		// Inject hook error
		case <-time.After(1 * time.Second):
			assert.Fail(t, "should have finished")
		}
	})

	t.Run("should fail and retry if HeaderByNumber fails", func(t *testing.T) {
		cancellableCtx, cancel := context.WithCancel(ctx)
		chain := newFakeChain()
		session := NewSession(chain, mockEthClient, mockClient, mockHook, mockOffsetManager, mockMetrics)
		bckoff := &backoffmock.MockIntervalBackoff{}
		session.bckOff = bckoff

		mockOffsetManager.EXPECT().GetLastBlockNumber(gomock.Any(), chain).Return(blockPosition.Uint64(), nil).AnyTimes()
		mockEthClient.EXPECT().EEAPrivPrecompiledContractAddr(gomock.Any(), chain.URL).Return(common.HexToAddress(eeaPrivPrecompiledContractAddr), nil).AnyTimes()
		mockEthClient.EXPECT().HeaderByNumber(gomock.Any(), chain.URL, nil).Return(nil, fmt.Errorf("error")).AnyTimes()

		// Start session
		exitErr := make(chan error)
		go func() {
			exitErr <- session.Run(cancellableCtx)
		}()
		go func() {
			<-time.After(200 * time.Microsecond)
			cancel()
		}()

		select {
		case <-exitErr:
			assert.True(t, bckoff.HasRetried())
		// Inject hook error
		case <-time.After(1 * time.Second):
			assert.Fail(t, "should have finished")
		}
	})

	t.Run("should not fail if BlockByNumber fails", func(t *testing.T) {
		cancellableCtx, cancel := context.WithCancel(ctx)
		chain := newFakeChain()
		session := NewSession(chain, mockEthClient, mockClient, mockHook, mockOffsetManager, mockMetrics)
		bckoff := &backoffmock.MockIntervalBackoff{}
		session.bckOff = bckoff

		mockOffsetManager.EXPECT().GetLastBlockNumber(gomock.Any(), chain).Return(blockPosition.Uint64(), nil).AnyTimes()
		mockEthClient.EXPECT().EEAPrivPrecompiledContractAddr(gomock.Any(), chain.URL).Return(common.HexToAddress(eeaPrivPrecompiledContractAddr), nil).AnyTimes()
		mockEthClient.EXPECT().HeaderByNumber(gomock.Any(), chain.URL, nil).Return(&types.Header{
			Number: newBlockPosition,
		}, nil).AnyTimes()
		mockEthClient.EXPECT().BlockByNumber(gomock.Any(), chain.URL, newBlockPosition).Return(nil, fmt.Errorf("error")).AnyTimes()

		// Start session
		exitErr := make(chan error)
		go func() {
			exitErr <- session.Run(cancellableCtx)
		}()
		go func() {
			<-time.After(200 * time.Microsecond)
			cancel()
		}()

		select {
		case <-exitErr:
			// Success if we have an error
		case <-time.After(1 * time.Second):
			assert.Fail(t, "should have finished")
		}
	})

	t.Run("should not fail if SearchJob fails", func(t *testing.T) {
		cancellableCtx, cancel := context.WithCancel(ctx)
		chain := newFakeChain()
		block := newFakeBlock(newBlockPosition, eeaPrivPrecompiledContractAddr)
		session := NewSession(chain, mockEthClient, mockClient, mockHook, mockOffsetManager, mockMetrics)
		bckoff := &backoffmock.MockIntervalBackoff{}
		session.bckOff = bckoff

		mockOffsetManager.EXPECT().GetLastBlockNumber(gomock.Any(), chain).Return(blockPosition.Uint64(), nil).AnyTimes()
		mockEthClient.EXPECT().EEAPrivPrecompiledContractAddr(gomock.Any(), chain.URL).Return(common.HexToAddress(eeaPrivPrecompiledContractAddr), nil).AnyTimes()
		mockEthClient.EXPECT().HeaderByNumber(gomock.Any(), chain.URL, nil).Return(&types.Header{
			Number: newBlockPosition,
		}, nil).AnyTimes()
		mockEthClient.EXPECT().BlockByNumber(gomock.Any(), chain.URL, newBlockPosition).Return(block, nil).AnyTimes()
		mockClient.EXPECT().
			SearchJob(gomock.Any(), &entities.JobFilters{
				TxHashes:  []string{txHashPrivate},
				ChainUUID: chain.UUID,
				Status:    entities.StatusPending,
			}).
			Return(nil, fmt.Errorf("search job error")).AnyTimes()

		// Start session
		exitErr := make(chan error)
		go func() {
			exitErr <- session.Run(cancellableCtx)
		}()
		go func() {
			<-time.After(200 * time.Microsecond)
			cancel()
		}()

		select {
		case <-exitErr:
			// Success if we have an error
		case <-time.After(1 * time.Second):
			assert.Fail(t, "should have finished")
		}
	})

	t.Run("should not fail if PrivateTransactionReceipt fails", func(t *testing.T) {
		cancellableCtx, cancel := context.WithCancel(ctx)
		chain := newFakeChain()
		chain.Listener.ExternalTxEnabled = true
		block := newFakeBlock(newBlockPosition, eeaPrivPrecompiledContractAddr)
		session := NewSession(chain, mockEthClient, mockClient, mockHook, mockOffsetManager, mockMetrics)
		bckoff := &backoffmock.MockIntervalBackoff{}
		session.bckOff = bckoff

		mockOffsetManager.EXPECT().GetLastBlockNumber(gomock.Any(), chain).Return(blockPosition.Uint64(), nil).AnyTimes()
		mockEthClient.EXPECT().EEAPrivPrecompiledContractAddr(gomock.Any(), chain.URL).Return(common.HexToAddress(eeaPrivPrecompiledContractAddr), nil).AnyTimes()
		mockEthClient.EXPECT().HeaderByNumber(gomock.Any(), chain.URL, nil).Return(&types.Header{
			Number: newBlockPosition,
		}, nil).AnyTimes()
		mockEthClient.EXPECT().BlockByNumber(gomock.Any(), chain.URL, newBlockPosition).Return(block, nil).AnyTimes()
		mockClient.EXPECT().
			SearchJob(gomock.Any(), &entities.JobFilters{
				TxHashes:  []string{txHashPrivate},
				ChainUUID: chain.UUID,
				Status:    entities.StatusPending,
			}).
			Return([]*txschedulertypes.JobResponse{}, nil).AnyTimes()
		mockEthClient.EXPECT().
			PrivateTransactionReceipt(gomock.Any(), chain.URL, common.HexToHash(txHashPrivate)).
			Return(nil, fmt.Errorf("private receipt error")).AnyTimes()

		// Start session
		exitErr := make(chan error)
		go func() {
			exitErr <- session.Run(cancellableCtx)
		}()
		go func() {
			<-time.After(200 * time.Microsecond)
			cancel()
		}()

		select {
		case <-exitErr:
			// Success if we have an error
		case <-time.After(1 * time.Second):
			assert.Fail(t, "should have finished")
		}
	})

	t.Run("should not fail if TransactionReceipt fails", func(t *testing.T) {
		cancellableCtx, cancel := context.WithCancel(ctx)
		chain := newFakeChain()
		chain.Listener.ExternalTxEnabled = true
		block := newFakeBlock(newBlockPosition, toAddress)
		session := NewSession(chain, mockEthClient, mockClient, mockHook, mockOffsetManager, mockMetrics)
		bckoff := &backoffmock.MockIntervalBackoff{}
		session.bckOff = bckoff

		mockOffsetManager.EXPECT().GetLastBlockNumber(gomock.Any(), chain).Return(blockPosition.Uint64(), nil).AnyTimes()
		mockEthClient.EXPECT().EEAPrivPrecompiledContractAddr(gomock.Any(), chain.URL).Return(common.HexToAddress(eeaPrivPrecompiledContractAddr), nil).AnyTimes()
		mockEthClient.EXPECT().HeaderByNumber(gomock.Any(), chain.URL, nil).Return(&types.Header{
			Number: newBlockPosition,
		}, nil).AnyTimes()
		mockEthClient.EXPECT().BlockByNumber(gomock.Any(), chain.URL, newBlockPosition).Return(block, nil).AnyTimes()
		mockClient.EXPECT().
			SearchJob(gomock.Any(), &entities.JobFilters{
				TxHashes:  []string{txHash},
				ChainUUID: chain.UUID,
				Status:    entities.StatusPending,
			}).
			Return([]*txschedulertypes.JobResponse{}, nil).AnyTimes()
		mockEthClient.EXPECT().
			TransactionReceipt(gomock.Any(), chain.URL, common.HexToHash(txHash)).
			Return(nil, fmt.Errorf("receipt error")).AnyTimes()

		// Start session
		exitErr := make(chan error)
		go func() {
			exitErr <- session.Run(cancellableCtx)
		}()
		go func() {
			<-time.After(200 * time.Microsecond)
			cancel()
		}()

		select {
		case <-exitErr:
			// Success if we have an error
		case <-time.After(1 * time.Second):
			assert.Fail(t, "should have finished")
		}
	})
}

func newFakeChain() *dynamic.Chain {
	backoff, _ := time.ParseDuration("1s")

	return &dynamic.Chain{
		UUID:     "chainUUID",
		TenantID: "tenantID",
		Name:     "chainName",
		URL:      "chainURL",
		ChainID:  "888",
		Listener: dynamic.Listener{
			StartingBlock:     0,
			CurrentBlock:      0,
			Depth:             0,
			Backoff:           backoff,
			ExternalTxEnabled: false,
		},
		Active: true,
	}
}

func newFakeBlock(blockPosition *big.Int, to string) *types.Block {
	txs := []*types.Transaction{types.NewTransaction(
		0,
		common.HexToAddress(to),
		big.NewInt(10000),
		uint64(21000),
		big.NewInt(20000),
		[]byte{},
	)}

	return types.NewBlock(&types.Header{Number: blockPosition}, txs, []*types.Header{}, []*types.Receipt{})
}

func newFakeReceipt() *ethereum.Receipt {
	return &ethereum.Receipt{
		TxHash:    "0xtxHash",
		BlockHash: "0xblockHash",
	}
}
