package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type mockEthClient struct {
	t                       *testing.T
	index                   int
	expected                []string
	mux                     *sync.Mutex
	callsBlockByNumber      []*big.Int
	callsTransactionReceipt []common.Hash
}

var (
	blockEnc = common.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0")
	block    types.Block
	_        = rlp.DecodeBytes(blockEnc, &block)
)

func (ec *mockEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	ec.mux.Lock()
	defer ec.mux.Unlock()
	ec.callsBlockByNumber = append(ec.callsBlockByNumber, number)
	defer func() { ec.index++ }()
	if ec.index < len(ec.expected) {
		if ec.expected[ec.index] == "mined" {
			return &block, nil
		}
		if ec.expected[ec.index] == "pending" {
			return nil, nil
		}
	}

	return nil, fmt.Errorf("Block missing")
}

func (ec *mockEthClient) TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	// Simulate some io time
	time.Sleep(time.Duration(50) * time.Millisecond)
	ec.mux.Lock()
	defer ec.mux.Unlock()
	ec.callsTransactionReceipt = append(ec.callsTransactionReceipt, hash)

	if block.Transaction(hash) != nil {
		return &types.Receipt{
			TxHash: hash,
		}, nil
	}

	return nil, fmt.Errorf("Receipt missing")
}

func TestBlockCursor(t *testing.T) {

	// Create a cursor
	cursor := &blockCursor{
		ec: &mockEthClient{
			t:        t,
			mux:      &sync.Mutex{},
			index:    0,
			expected: []string{"mined"},
		},
		next: big.NewInt(0),
	}

	block, err := cursor.Next()
	if err != nil {
		t.Errorf("Next: expected to get block but got %v", err)
	}

	if !reflect.DeepEqual(block.Hash(), common.HexToHash("80bd976d96ef1da0251150e741fd596d3e580be70b02a4757554a452c17edfe1")) {
		t.Errorf("Next: Hash mismatch got %v, want %v", block.Hash().Hex(), "0x80bd976d96ef1da0251150e741fd596d3e580be70b02a4757554a452c17edfe1")
	}

	if cursor.next.Uint64() != 1 {
		t.Errorf("Next: Cursor should have been incremented")
	}

	block, err = cursor.Next()
	if err == nil || block != nil {
		t.Errorf("Next: expected an error")
	}

	if cursor.next.Uint64() != 1 {
		t.Errorf("Next: Cursor should not have been incremented")
	}

	cursor.Set(big.NewInt(10))
	if cursor.next.Uint64() != 10 {
		t.Errorf("Next: Cursor should have been set to %v but got %v", 10, cursor.next.Uint64())
	}
}

func TestBlockConsumer(t *testing.T) {
	// Create mock client
	ec := &mockEthClient{
		t:        t,
		mux:      &sync.Mutex{},
		index:    0,
		expected: []string{"mined", "mined", "pending", "pending", "mined", "pending", "mined", "mined", "pending", "mined"},
	}

	// Create cursor & block consumer
	cursor := &blockCursor{
		ec:   ec,
		next: big.NewInt(0),
	}
	config := &TxListenerConfig{}
	config.Block.Delay = 10 * time.Millisecond
	bc := newBlockConsumer(cursor, config)

	// Run feeder
	go bc.feeder()

	blocks := []*types.Block{}
	for block := range bc.Blocks() {
		blocks = append(blocks, block)
	}

	if len(blocks) != 6 {
		t.Errorf("BlockConsumer: expected 6 blocks but got %v", len(blocks))
	}

	if len(ec.callsBlockByNumber) != 11 {
		t.Errorf("BlockConsumer: expected %v calls on BlockByNumber but got %v", 11, len(ec.callsBlockByNumber))
	}
}

func TestBlockConsumerInterupted(t *testing.T) {
	// Create mock client
	ec := &mockEthClient{
		t:        t,
		mux:      &sync.Mutex{},
		index:    0,
		expected: []string{"mined", "mined", "pending", "pending", "mined", "pending", "mined", "mined", "pending", "mined"},
	}

	// Create cursor & block consumer
	cursor := &blockCursor{
		ec:   ec,
		next: big.NewInt(0),
	}
	config := &TxListenerConfig{}
	config.Block.Delay = 10 * time.Millisecond
	bc := newBlockConsumer(cursor, config)

	// Run feeder
	go bc.feeder()

	go func() {
		// Interupt block consumer
		time.Sleep(25 * time.Millisecond)
		bc.Close()
	}()

	blocks := []*types.Block{}
	for block := range bc.Blocks() {
		blocks = append(blocks, block)
	}

	if len(blocks) != 3 {
		t.Errorf("BlockConsumer: expected %v blocks but got %v", 3, len(blocks))
	}

	if len(ec.callsBlockByNumber) != 6 {
		t.Errorf("BlockConsumer: expected %v calls on BlockByNumber but got %v", 6, len(ec.callsBlockByNumber))
	}
}

func TestTxListener(t *testing.T) {
	// Create mock client
	ec := &mockEthClient{
		t:        t,
		mux:      &sync.Mutex{},
		index:    0,
		expected: []string{"mined", "mined", "pending", "pending", "mined", "pending", "mined", "mined", "pending", "mined"},
	}

	config := &TxListenerConfig{}
	config.Block.Delay = 50 * time.Millisecond
	config.Receipts.Count = 200

	// Create txListener
	l := NewTxListener(ec, config)

	receipts := []*types.Receipt{}
	for receipt := range l.Receipts() {
		receipts = append(receipts, receipt)
	}

	if len(receipts) != len(block.Transactions())*6 {
		t.Errorf("TxListener: expected %v receipts but got %v", 6, len(receipts))
	}

	// Ensure Receipts have processed in expected order
	for i := 0; i < 6; i++ {
		for j, tx := range block.Transactions() {
			if tx.Hash().Hex() != receipts[len(block.Transactions())*i+j].TxHash.Hex() {
				t.Errorf("TxListener: expected TxHash %v but got %v", tx.Hash().Hex(), receipts[len(block.Transactions())*i+j].TxHash.Hex())
			}
		}
	}
}
