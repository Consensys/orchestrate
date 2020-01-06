package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
	offset "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset/memory"
)

type testKey string

const (
	ecError = testKey("ec.error")
)

type EthClientV2 struct {
	mux *sync.RWMutex

	blocks    []*ethtypes.Block
	blocksIdx map[string]int

	txs map[string]*ethtypes.Transaction

	errors chan error

	chainID *big.Int
}

func NewEthClientV2(chainID *big.Int) *EthClientV2 {
	return &EthClientV2{
		mux:       &sync.RWMutex{},
		blocksIdx: make(map[string]int),
		txs:       make(map[string]*ethtypes.Transaction),
		errors:    make(chan error, 1),
		chainID:   chainID,
	}
}

func (ec *EthClientV2) Mine(block *ethtypes.Block) {
	ec.mux.Lock()
	defer ec.mux.Unlock()

	// Update block number
	header := block.Header()
	header.Number = big.NewInt(int64(len(ec.blocks)))
	b := ethtypes.NewBlock(header, block.Transactions(), block.Uncles(), nil)

	ec.blocks = append(ec.blocks, b)
	ec.blocksIdx[b.Hash().Hex()] = len(ec.blocks)

	for _, tx := range b.Transactions() {
		ec.txs[tx.Hash().Hex()] = tx
	}
}

func (ec *EthClientV2) Errors() <-chan error {
	return ec.errors
}

func (ec *EthClientV2) getError(ctx context.Context) error {
	err, ok := ctx.Value(ecError).(error)
	if ok {
		return err
	}

	select {
	case err := <-ec.errors:
		return err
	default:
		return nil
	}
}

func (ec *EthClientV2) BlockByHash(ctx context.Context, _ string, hash ethcommon.Hash) (*ethtypes.Block, error) {
	if err := ec.getError(ctx); err != nil {
		return nil, err
	}

	ec.mux.RLock()
	defer ec.mux.RUnlock()

	if idx, ok := ec.blocksIdx[hash.Hex()]; ok {
		return ec.blocks[idx], nil
	}

	return nil, errors.NotFoundError("block not found")
}

func (ec *EthClientV2) BlockByNumber(ctx context.Context, _ string, number *big.Int) (*ethtypes.Block, error) {
	if err := ec.getError(ctx); err != nil {
		return nil, err
	}

	ec.mux.RLock()
	defer ec.mux.RUnlock()

	if number == nil {
		return ec.blocks[len(ec.blocks)-1], nil
	}

	if idx := int(number.Uint64()); idx < len(ec.blocks) {
		return ec.blocks[idx], nil
	}

	return nil, errors.NotFoundError("block not found")
}

func (ec *EthClientV2) HeaderByHash(ctx context.Context, url string, hash ethcommon.Hash) (*ethtypes.Header, error) {
	block, err := ec.BlockByHash(ctx, url, hash)
	if err != nil {
		return nil, err
	}
	return block.Header(), nil
}

func (ec *EthClientV2) HeaderByNumber(ctx context.Context, url string, number *big.Int) (*ethtypes.Header, error) {
	block, err := ec.BlockByNumber(ctx, url, number)
	if err != nil {
		return nil, err
	}
	return block.Header(), nil
}

func (ec *EthClientV2) TransactionByHash(ctx context.Context, _ string, hash ethcommon.Hash) (tx *ethtypes.Transaction, isPending bool, err error) {
	if err := ec.getError(ctx); err != nil {
		return nil, false, err
	}

	ec.mux.RLock()
	defer ec.mux.RUnlock()

	if tx, ok := ec.txs[hash.Hex()]; ok {
		return tx, false, nil
	}

	return nil, false, errors.NotFoundError("tx not found")
}

func (ec *EthClientV2) TransactionReceipt(ctx context.Context, _ string, txHash ethcommon.Hash) (*ethtypes.Receipt, error) {
	if err := ec.getError(ctx); err != nil {
		return nil, err
	}

	ec.mux.RLock()
	defer ec.mux.RUnlock()

	if tx, ok := ec.txs[txHash.Hex()]; ok {
		return &ethtypes.Receipt{TxHash: tx.Hash()}, nil
	}

	return nil, errors.NotFoundError("receipt not found")
}

func (ec *EthClientV2) Network(ctx context.Context, _ string) (*big.Int, error) {
	if err := ec.getError(ctx); err != nil {
		return nil, err
	}
	return ec.chainID, nil
}

func (ec *EthClientV2) SyncProgress(ctx context.Context, _ string) (*eth.SyncProgress, error) {
	if err := ec.getError(ctx); err != nil {
		return nil, err
	}
	return &eth.SyncProgress{
		StartingBlock: 0,
		CurrentBlock:  uint64(len(ec.blocks)),
		HighestBlock:  uint64(len(ec.blocks)),
	}, nil
}

type hookCall struct {
	node     *dynamic.Node
	block    *ethtypes.Block
	receipts []*ethtypes.Receipt
}

type MockHook struct {
	Calls  chan *hookCall
	Errors chan error
}

func NewMockHook() *MockHook {
	return &MockHook{
		Calls:  make(chan *hookCall, 10),
		Errors: make(chan error, 1),
	}
}

func (hk *MockHook) AfterNewBlock(_ context.Context, node *dynamic.Node, block *ethtypes.Block, receipts []*ethtypes.Receipt) error {
	hk.Calls <- &hookCall{
		node:     node,
		block:    block,
		receipts: receipts,
	}

	select {
	case err := <-hk.Errors:
		return err
	default:
		return nil
	}
}

func (hk *MockHook) getCall(timeout time.Duration) (*hookCall, error) {
	select {
	case call := <-hk.Calls:
		return call, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("hook timeout")
	}
}

func TestGetChainTip(t *testing.T) {
	ec := NewEthClientV2(big.NewInt(1))

	// Mine 2 blocks
	blockEnc := ethcommon.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0")
	var block ethtypes.Block
	_ = rlp.DecodeBytes(blockEnc, &block)
	ec.Mine(&block)
	ec.Mine(&block)

	sess := &Session{
		ec:   ec,
		Node: &dynamic.Node{Listener: &dynamic.Listener{}},
	}

	tip, err := sess.getChainTip(context.Background())
	assert.NoError(t, err, "getChainTip should not error")
	assert.Equal(t, uint64(1), tip, "Tip should be correct")
}

func TestFetchReceipt(t *testing.T) {
	ec := NewEthClientV2(big.NewInt(1))

	// Mine 1 block
	blockEnc := ethcommon.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0")
	var block ethtypes.Block
	_ = rlp.DecodeBytes(blockEnc, &block)
	ec.Mine(&block)

	sess := &Session{
		ec:   ec,
		Node: &dynamic.Node{},
	}

	// Unkwnon transaction
	future := sess.fetchReceipt(context.Background(), ethcommon.Hash{})
	select {
	case <-future.Err():
	case <-future.Result():
		t.Errorf("Future should have errored")
	}
	future.Close()

	// Know receipt
	future = sess.fetchReceipt(
		context.Background(),
		ethcommon.HexToHash("0xfbf12011cab2a6c12e1ee895495f2d1aa534b2dc8abcfc10fff88356e5b990fa"),
	)
	select {
	case err := <-future.Err():
		t.Errorf("Future should not error but got %v", err)
	case res := <-future.Result():
		assert.Equal(t, "0xfbf12011cab2a6c12e1ee895495f2d1aa534b2dc8abcfc10fff88356e5b990fa", res.(*ethtypes.Receipt).TxHash.Hex(), "Receipt hash should be correct")
	}
	future.Close()
}

func TestFetchBlock(t *testing.T) {
	ec := NewEthClientV2(big.NewInt(1))

	// Mine 1 block
	blockEnc := ethcommon.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0")
	var block ethtypes.Block
	_ = rlp.DecodeBytes(blockEnc, &block)
	ec.Mine(&block)

	sess := &Session{
		ec:   ec,
		Node: &dynamic.Node{},
	}

	future := sess.fetchBlock(context.Background(), 0)
	select {
	case err := <-future.Err():
		t.Errorf("Future should not error but got %v", err)
	case res := <-future.Result():
		block := res.(*fetchedBlock)
		assert.NotNil(t, block, "Result block should not be nil")
		assert.Equal(t, "0xff4f5cd9a03569e8e6d32af4726d1b9ea1a248f69a04307f76896a24fe7be09d", block.block.Hash().Hex(), "Block hash should be correct(")
		assert.Len(t, block.receipts, 9, "Receipts should have been fetched properly")
	}
	future.Close()
}

func TestInit(t *testing.T) {
	ec := NewEthClientV2(big.NewInt(1))

	// Mine 2 block
	blockEnc := ethcommon.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0")
	var block ethtypes.Block
	_ = rlp.DecodeBytes(blockEnc, &block)
	ec.Mine(&block)
	ec.Mine(&block)

	hk := NewMockHook()
	offsets := offset.NewManager()
	builder := NewSessionBuilder(hk, offsets, ec)

	// Test 1: init from latest
	sess := builder.newSession(
		&dynamic.Node{
			Listener: &dynamic.Listener{
				BlockPosition: -1,
			},
		},
	)

	err := sess.init(context.Background())
	assert.NoError(t, err, "Init should not error")
	assert.Equal(t, uint64(1), sess.currentChainTip, "#1 Chain tip should be correct")
	assert.Equal(t, uint64(1), sess.blockPosition, "#1 blockPosition should be correct")
	assert.Equal(t, uint64(1), sess.Node.ChainID.Uint64(), "#1 node chain ID should have been set")

	// Test 2: init from block 2
	sess = builder.newSession(
		&dynamic.Node{
			Listener: &dynamic.Listener{
				BlockPosition: 2,
			},
		},
	)

	err = sess.init(context.Background())
	assert.NoError(t, err, "Init should not error")
	assert.Equal(t, uint64(0), sess.currentChainTip, "#2 Chain tip should be correct")
	assert.Equal(t, uint64(2), sess.blockPosition, "#2 blockPosition should be correct")

	// Test 3: init with offset manager initialized
	_ = offsets.SetLastBlockNumber(context.Background(), &dynamic.Node{ID: "test-init"}, 14)
	sess = builder.newSession(
		&dynamic.Node{
			ID: "test-init",
			Listener: &dynamic.Listener{
				BlockPosition: 2,
			},
		},
	)

	err = sess.init(context.Background())
	assert.NoError(t, err, "Init should not error")
	assert.Equal(t, uint64(0), sess.currentChainTip, "#2 Chain tip should be correct")
	assert.Equal(t, uint64(15), sess.blockPosition, "#2 blockPosition should be correct")
}

func TestRun(t *testing.T) {
	hk := NewMockHook()
	offsets := offset.NewManager()

	// Initialize MockLedgerReader with 1st mined block
	ec := NewEthClientV2(big.NewInt(1))
	blockEnc := ethcommon.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0")
	var block ethtypes.Block
	_ = rlp.DecodeBytes(blockEnc, &block)
	ec.Mine(&block)

	// Create builder
	builder := NewSessionBuilder(hk, offsets, ec)

	// New session starting on block one
	node := &dynamic.Node{
		ID: "test-node",
		Listener: &dynamic.Listener{
			BlockPosition: 1,
			Backoff:       10 * time.Millisecond,
		},
	}
	sess := builder.newSession(node)

	// Start session
	errChan := make(chan error)
	cancelableCtx, cancel := context.WithCancel(context.Background())
	go func() {
		errChan <- sess.Run(cancelableCtx)
	}()

	// Make sure hook has been properly called
	_, err := hk.getCall(100 * time.Millisecond)
	assert.Error(t, err, "#1 Hook should not have been called yet")

	// Mine new block then wait for hook to be called
	ec.Mine(&block)
	call, err := hk.getCall(100 * time.Millisecond)
	assert.NoError(t, err, "#2 Hook should not have been called yet")
	assert.Equal(t, "0xbb198635820f9fff08f1af3aea743001a2469a6fb8b1cb0881995a8ea7a26b32", call.block.Hash().Hex(), "#2 Block hash should be correct")
	assert.Len(t, call.receipts, 9, "#2 Receipt count should be correct")

	// Mine new block then wait for hook to be called
	ec.Mine(&block)
	call, err = hk.getCall(100 * time.Millisecond)
	assert.NoError(t, err, "#3 Hook should not have been called yet")
	assert.Equal(t, "0x8f57b12395c3beb3382f1c2acc585566774e228d87bea349798f5793d01150af", call.block.Hash().Hex(), "#3 Block hash should be correct")
	assert.Len(t, call.receipts, 9, "#3 Receipt count should be correct")

	// Cancel context
	cancel()

	select {
	case err := <-errChan:
		assert.NoError(t, err, "Session should not have error")
	case <-time.After(time.Second):
		t.Errorf("TestRun: session should have completed")
	}
	close(errChan)
	lastBlock, _ := offsets.GetLastBlockNumber(context.Background(), node)
	assert.Equal(t, int64(2), lastBlock, "Offset manager should have properly updated block processed")
}

func TestRunWithError(t *testing.T) {
	hk := NewMockHook()
	offsets := offset.NewManager()

	// Initialize MockLedgerReader with 1st mined block
	ec := NewEthClientV2(big.NewInt(1))
	blockEnc := ethcommon.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0")
	var block ethtypes.Block
	_ = rlp.DecodeBytes(blockEnc, &block)
	ec.Mine(&block)

	// Create builder
	builder := NewSessionBuilder(hk, offsets, ec)

	// New session starting on block one
	node := &dynamic.Node{
		ID: "test-node",
		Listener: &dynamic.Listener{
			BlockPosition: 1,
			Backoff:       10 * time.Millisecond,
		},
	}
	sess := builder.newSession(node)

	// Start session
	errChan := make(chan error)
	go func() {
		errChan <- sess.Run(context.Background())
	}()

	// Inject hook error
	hk.Errors <- fmt.Errorf("hook-error")

	// Mine new block then wait for hook to be called
	ec.Mine(&block)

	select {
	case err := <-errChan:
		assert.Error(t, err, "Session should have errored")
	case <-time.After(time.Second):
		t.Errorf("TestRunError: session should have completed")
	}
}
