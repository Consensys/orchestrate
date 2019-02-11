package listener

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type MockEthClient struct {
	blocks []*types.Block

	mux  *sync.RWMutex
	head uint64
}

func NewMockEthClient(blocks []*types.Block) *MockEthClient {
	return &MockEthClient{
		blocks: blocks,
		mux:    &sync.RWMutex{},
	}
}

func (ec *MockEthClient) mine() {
	ec.mux.Lock()
	defer ec.mux.Unlock()

	if int(ec.head) < len(ec.blocks) {
		ec.head++
	}
}

type MockKey string

func (ec *MockEthClient) BlockByNumber(ctx context.Context, chainID *big.Int, number *big.Int) (*types.Block, error) {
	_, ok := ctx.Value(MockKey("error")).(error)
	if ok {
		return nil, fmt.Errorf("MockEthClient: Error on BlockByNumber")
	}

	// Simulate io time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(2 * time.Millisecond):
		ec.mux.RLock()
		defer ec.mux.RUnlock()

		if number == nil {
			number = big.NewInt(int64(ec.head))
		}

		if number.Uint64() <= ec.head {
			block := ec.blocks[number.Uint64()]
			header := types.CopyHeader(block.Header())
			header.Number = number
			blck := types.NewBlockWithHeader(header)
			return blck.WithBody(block.Transactions(), block.Uncles()), nil
		}

		if number.Uint64() > ec.head {
			return nil, nil
		}
		return nil, fmt.Errorf("Error")
	}
}

func (ec *MockEthClient) HeaderByNumber(ctx context.Context, chainID *big.Int, number *big.Int) (*types.Header, error) {
	_, ok := ctx.Value(MockKey("error")).(error)
	if ok {
		return nil, fmt.Errorf("MockEthClient: Error on SyncProgress")
	}

	// Simulate io time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(2 * time.Millisecond):
		ec.mux.RLock()
		defer ec.mux.RUnlock()
		if number == nil {
			number = big.NewInt(int64(ec.head))
		}

		if number.Uint64() <= ec.head {
			block := ec.blocks[number.Uint64()]
			header := types.CopyHeader(block.Header())
			header.Number = number
			return header, nil
		}

		if number.Uint64() > ec.head {
			return nil, nil
		}

		return nil, fmt.Errorf("Error")
	}
}

func (ec *MockEthClient) TransactionReceipt(ctx context.Context, chainID *big.Int, txHash common.Hash) (*types.Receipt, error) {
	_, ok := ctx.Value(MockKey("error")).(error)
	if ok {
		return nil, fmt.Errorf("MockEthClient: Error on TransactionReceipt")
	}

	// Simulate io time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(2 * time.Millisecond):
		ec.mux.RLock()
		defer ec.mux.RUnlock()
		for _, block := range ec.blocks[:ec.head+1] {
			if block.Transaction(txHash) != nil {
				return &types.Receipt{
					TxHash: txHash,
				}, nil
			}
		}
		return nil, nil
	}
}

// TODO: update with disctinct blocks
var blocksEnc = [][]byte{
	common.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0"),
	common.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0"),
	common.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0"),
	common.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0"),
}

func TestMockEthClient(t *testing.T) {
	blocks := []*types.Block{}
	for _, blockEnc := range blocksEnc {
		var block types.Block
		rlp.DecodeBytes(blockEnc, &block)
		blocks = append(blocks, &block)
	}
	mec := NewMockEthClient(blocks)

	b, err := mec.BlockByNumber(context.Background(), big.NewInt(1), big.NewInt(0))
	if b == nil || err != nil {
		t.Errorf("MockEthClient #1: Got %v %v", b, err)
	}

	h, err := mec.HeaderByNumber(context.Background(), big.NewInt(1), big.NewInt(0))
	if h.Number.Int64() != 0 || err != nil {
		t.Errorf("MockEthClient #1: Head at %v", h.Number.Int64())
	}

	b, err = mec.BlockByNumber(context.Background(), big.NewInt(1), big.NewInt(1))
	if b != nil || err != nil {
		t.Errorf("MockEthClient #2: Got %v %v", b, err)
	}

	mec.mine()

	b, err = mec.BlockByNumber(context.Background(), big.NewInt(1), big.NewInt(1))
	if b == nil || err != nil {
		t.Errorf("MockEthClient #3: Got %v %v", b, err)
	}

	h, err = mec.HeaderByNumber(context.Background(), big.NewInt(1), nil)
	if h.Number.Int64() != 1 || err != nil {
		t.Errorf("MockEthClient #3: Head at %v", h.Number.Int64())
	}

	ctx := context.WithValue(context.Background(), MockKey("error"), fmt.Errorf("Error"))
	b, err = mec.BlockByNumber(ctx, big.NewInt(1), big.NewInt(1))
	if b != nil || err == nil {
		t.Errorf("MockEthClient #4: Got %v %v", b, err)
	}
}

func TestBaseTracker(t *testing.T) {
	blocks := []*types.Block{}
	for _, blockEnc := range blocksEnc {
		var block types.Block
		rlp.DecodeBytes(blockEnc, &block)
		blocks = append(blocks, &block)
	}
	mec := NewMockEthClient(blocks)

	tracker := BaseTracker{
		ec:      mec,
		chainID: big.NewInt(1),
		depth:   2,
	}

	if tracker.ChainID().Text(16) != "1" {
		t.Errorf("BaseTracker: expected chainID %q but got %q", "1", tracker.ChainID().Text(16))
	}

	head, _ := tracker.HighestBlock(context.Background())
	if head != 0 {
		t.Errorf("BaseTracker #1: Head at %v", head)
	}

	mec.mine()
	mec.mine()
	mec.mine()

	head, _ = tracker.HighestBlock(context.Background())
	if head != 1 {
		t.Errorf("BaseTracker #2: Head at %v", head)
	}

	tracker = BaseTracker{
		ec:      mec,
		chainID: big.NewInt(1),
		depth:   0,
	}

	head, _ = tracker.HighestBlock(context.Background())
	if head != 3 {
		t.Errorf("BaseTracker #3: Head at %v", head)
	}

	ctx := context.WithValue(context.Background(), MockKey("error"), fmt.Errorf("Error"))
	_, err := tracker.HighestBlock(ctx)
	if err == nil {
		t.Errorf("BaseTracker #4: Expected an error")
	}
}

func TestBlockCursorFetchReceipt(t *testing.T) {
	blocks := []*types.Block{}
	for _, blockEnc := range blocksEnc {
		var block types.Block
		rlp.DecodeBytes(blockEnc, &block)
		blocks = append(blocks, &block)
	}
	mec := NewMockEthClient(blocks)

	tracker := &BaseTracker{
		ec:      mec,
		chainID: big.NewInt(1),
		depth:   0,
	}

	// Initialize cursor
	config := NewConfig()
	bc := newBlockCursorFromTracker(mec, tracker, -1, config)

	future := bc.fetchReceipt(context.Background(), common.HexToHash("0x8305d6f07eaced88f5f8f52d5acceedb07568c6ca6c956bef461ed3d6e77686b"))

	receipt := (<-future.res).(*TxListenerReceipt)
	if receipt.ChainID.Text(10) != "1" {
		t.Errorf("GetReceipt: got ChainID %q", receipt.ChainID.Text(10))
	}

	// Force error
	ctx := context.WithValue(context.Background(), MockKey("error"), fmt.Errorf("Error"))
	future = bc.fetchReceipt(ctx, common.HexToHash("0x8305d6f07eaced88f5f8f52d5acceedb07568c6ca6c956bef461ed3d6e77686b"))
	err := <-future.err
	if err == nil {
		t.Errorf("GetReceipt: got Err %v", err)
	}

	// Error on unknwon receipt
	future = bc.fetchReceipt(context.Background(), common.HexToHash("0xbabebeef"))
	err = <-future.err
	if err == nil {
		t.Errorf("GetReceipt: got Err %v", err)
	}
}

func TestBlockCursorFetchBlock(t *testing.T) {
	blocks := []*types.Block{}
	for _, blockEnc := range blocksEnc {
		var block types.Block
		rlp.DecodeBytes(blockEnc, &block)
		blocks = append(blocks, &block)
	}
	mec := NewMockEthClient(blocks)

	tracker := &BaseTracker{
		ec:      mec,
		chainID: big.NewInt(1),
		depth:   0,
	}

	// Initialize cursor
	config := NewConfig()
	bc := newBlockCursorFromTracker(mec, tracker, -1, config)

	future := bc.fetchBlock(context.Background(), 4)

	err := <-future.err
	if err == nil {
		t.Errorf("GetBlock: got Err %v", err)
	}

	future = bc.fetchBlock(context.Background(), 0)
	block := (<-future.res).(*TxListenerBlock)
	if block.ChainID.Text(10) != "1" {
		t.Errorf("GetBlock: got ChainID %q", block.ChainID.Text(10))
	}

	if len(block.receipts) != 9 {
		t.Errorf("GetBlock: expected %v receipts but got %v", 9, block.receipts)
	}

	for i, receipt := range block.receipts {
		if receipt.BlockHash.Hex() != block.Hash().Hex() {
			t.Errorf("GetBlock: expected blockhash %v but got %v", block.Hash().Hex(), receipt.BlockHash.Hex())
		}

		if receipt.BlockNumber != block.Number().Int64() {
			t.Errorf("GetBlock: expected BlockNumber %v but got %v", block.NumberU64(), receipt.BlockNumber)
		}

		if receipt.TxIndex != uint64(i) {
			t.Errorf("GetBlock: expected TxIndex %v but got %v", i, receipt.TxIndex)
		}
	}

	// Force error
	ctx := context.WithValue(context.Background(), MockKey("error"), fmt.Errorf("Error"))
	future = bc.fetchBlock(ctx, 0)

	err = <-future.err
	if err == nil {
		t.Errorf("GetBlock: got Err %v", err)
	}
}

func TestBlockCursorNext(t *testing.T) {
	blocks := []*types.Block{}
	for _, blockEnc := range blocksEnc {
		var block types.Block
		rlp.DecodeBytes(blockEnc, &block)
		blocks = append(blocks, &block)
	}
	mec := NewMockEthClient(blocks)

	tracker := &BaseTracker{
		ec:      mec,
		chainID: big.NewInt(1),
		depth:   0,
	}

	// Initialize cursor
	config := NewConfig()
	bc := newBlockCursorFromTracker(mec, tracker, -1, config)

	// No future have been so Next() should return false
	next := bc.Next(context.Background())
	if next {
		t.Errorf("Next #1: expected nothing ready")
	}

	// We feed a succesful future so Next() should return true
	future := &Future{
		res: make(chan interface{}),
		err: make(chan error),
	}
	go func() { future.res <- &TxListenerBlock{} }()
	bc.blockFutures <- future

	next = bc.Next(context.Background())
	if !next {
		t.Errorf("Next #2: expected a block ready")
	}

	// We feed a failing future so Next() should return true
	future = &Future{
		res: make(chan interface{}),
		err: make(chan error),
	}
	go func() { future.err <- fmt.Errorf("Test Error") }()
	bc.blockFutures <- future

	next = bc.Next(context.Background())
	if next {
		t.Errorf("Next #3: expected no block ready")
	}

	expected := "tx-listener: error while listening on chain 0x1: Test Error"
	if bc.err.Error() != expected {
		t.Errorf("Next #3: expected error %q bug got %q", expected, bc.err)
	}

	// We feed context that we cancel
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Cancel asynchronously
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	future = &Future{
		res: make(chan interface{}),
		err: make(chan error),
	}
	bc.blockFutures <- future

	next = bc.Next(ctx)
	if next {
		t.Errorf("Next #4: expected no block ready")
	}

	expected = "tx-listener: error while listening on chain 0x1: context canceled"
	if bc.err.Error() != expected {
		t.Errorf("Next #4: expected error %q bug got %q", expected, bc.err)
	}
}

func TestBlockCursor(t *testing.T) {
	blocks := []*types.Block{}
	for _, blockEnc := range blocksEnc {
		var block types.Block
		rlp.DecodeBytes(blockEnc, &block)
		blocks = append(blocks, &block)
	}
	mec := NewMockEthClient(blocks)

	// Initialize cursor
	config := NewConfig()
	config.BlockCursor.Backoff = 100 * time.Millisecond
	bc := NewBlockCursor(mec, big.NewInt(1), 0, config)
	defer bc.Close()

	if bc.ChainID().Text(16) != "1" {
		t.Errorf("BlockCursor expected chainID %q but got %q", "1", bc.ChainID().Text(16))
	}

	next := bc.Next(context.Background())
	if next {
		t.Errorf("BlockCursor #1: did not expected a block ready")
	}

	// Sleep waiting for Cursor to start
	time.Sleep(10 * time.Millisecond)

	next = bc.Next(context.Background())
	if !next {
		t.Errorf("BlockCursor #2: expected a ready block")
	}

	if len(bc.block.receipts) != 9 {
		t.Errorf("BlockCursor #3: expected 9 receipts but got %v", len(bc.block.receipts))
	}

	// No new mined block so next should not be ready
	next = bc.Next(context.Background())
	if next {
		t.Errorf("BlockCursor #4: did not expected a block ready")
	}

	// We simulate 2 new mined blocks
	mec.mine()
	mec.mine()

	// Sleep less than blockCursor backoff time, so Next should not be ready
	time.Sleep(50 * time.Millisecond)
	next = bc.Next(context.Background())
	if next {
		t.Errorf("BlockCursor #5: did not expected a block ready")
	}

	// Sleep again to pass block cursor backoff time, so Next should be ready
	time.Sleep(100 * time.Millisecond)
	next = bc.Next(context.Background())
	if !next {
		t.Errorf("BlockCursor #6: expected a ready block")
	}

	// A second block should be ready as we mined twice
	next = bc.Next(context.Background())
	if !next {
		t.Errorf("BlockCursor #7: expected a ready block")
	}
}
