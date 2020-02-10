package kafka

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	uuid "github.com/satori/go.uuid"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	"github.com/Shopify/sarama/mocks"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
	crc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry/client/mocks"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store/client/mocks"
)

type testKey string

type MockChainStateReader struct {
}

func (ec *MockChainStateReader) BalanceAt(ctx context.Context, url string, account ethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	return nil, nil
}

func (ec *MockChainStateReader) StorageAt(ctx context.Context, url string, account ethcommon.Address, key ethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	return nil, nil
}

func (ec *MockChainStateReader) CodeAt(ctx context.Context, url string, account ethcommon.Address, blockNumber *big.Int) ([]byte, error) {
	if err, ok := ctx.Value(testKey("codeAtError")).(error); ok {
		return nil, err
	}
	return ethcommon.Hex2Bytes("0xabcd"), nil
}

func (ec *MockChainStateReader) NonceAt(ctx context.Context, url string, account ethcommon.Address, blockNumber *big.Int) (uint64, error) {
	return 0, nil
}

func (ec *MockChainStateReader) PendingBalanceAt(ctx context.Context, url string, account ethcommon.Address) (*big.Int, error) {
	return nil, nil
}

func (ec *MockChainStateReader) PendingStorageAt(ctx context.Context, url string, account ethcommon.Address, key ethcommon.Hash) ([]byte, error) {
	return nil, nil
}

func (ec *MockChainStateReader) PendingCodeAt(ctx context.Context, url string, account ethcommon.Address) ([]byte, error) {
	return nil, nil
}

func (ec *MockChainStateReader) PendingNonceAt(ctx context.Context, url string, account ethcommon.Address) (uint64, error) {
	return 0, nil
}

func TestHook(t *testing.T) {
	// Initialize hook
	conf := &Config{
		OutTopic: "test-topic-decoder",
	}
	registry := crc.New()
	ec := &MockChainStateReader{}
	store := clientmock.New()
	producer := mocks.NewSyncProducer(t, nil)

	hk := NewHook(conf, registry, ec, store, producer)

	blockEnc := ethcommon.FromHex("0xf90600f90218a0e19f046955d37c5e23c2857cbeb602b72eeeb47b1539d604e88c16053480f41ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479480cb31b335dd587d1cc49723837f448c3c2e4736a02a30b9f172a3c58dbbaa5e890243e9d94fe669f50cbf237c34d41e8a3c150807a01e16eb6a5be337178a8b41b2dbc8481af9b4deb09dc25fb3e399c698e56ef560a04416bfc7541f873da23002d8b26a55f73e1dbba48c1d0b46bf366d055549b021b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085012590553683492c22837a12008302e248845c36859c99d883010900846765746888676f312e31312e34856c696e7578a0fff3f838abb411d1bfaa65a9a3d1e7c162d9e8293802c30a73ff0064d42af53f88ba5707f7725a3c0ef903e1f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba06693c6a8f27c38aa559ffd7952a3cc06330fa6f3b75f966f3b782acb5a12d629a04d74f460391f4e843134c524ca304d9d8b95fa4e72173e3e58316469a9d98ae6f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba0361a8dd7c6ba0583bd469fc2ad5e360ca185e66e0caa28329bffa41a26b128a9a02f7e0823a3e182dde15e9ef9e64da11ee55b2940f887221154b821c58f09cb80f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08fd9cca3c8b509da67e138062d67325e6986a12620ecfb77ef1bc09578c218a5a00c407b0be555900c97afbe3f2022e5a49fdf84dbc25c7a906b5678550b5593f0f86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ca091d4b169b328e82a9ac6a36fa4703865e76b66f668cf86b35a39aded9586455ba00d173820d80ca8b0b4b103e820e51e2c145f753e168d615526653ca022478f9ef86f840128c364843b9aca00825208945cc31379a0a1d1a56c7a35cfcdeb96ca83c95277880de0b6b3a7640000801ca0bef5dfcccf430b07ce9b0d89ff31b7ee765586b376991cb39478a65f622c7753a03549afb66bfeb7e31bbfb31b8510ded604e559e0e0811c38a5e90e0841180809f86f840128c365843b9aca008252089455efceae4188f18e39c4cebd1d0a1502706aebd9880de0b6b3a7640000801ba006f4f786295bc218c187f2ee1cff23470745d6b4efc6188a28eebbea3136d447a05d95382232701baf7b4636f8b5a7b43d53d0e60d9f2953fc1b44e975a3be7d7cf86f840128c366843b9aca00825208949244af76c192ec3e525e97557da454ce1fcfe914880de0b6b3a7640000801ba0c24e71aff9952f667481eaf613e64fa6e5a1d566fbc843e41619d3a99ea7edcba05920f45a7d669b555373a7ee064b60479182961c6a15d056a9a64e55b635bdccf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba006da2cb18cfd311bb2b149ad85cd3ebf02c3db7178fc97e06db32c0743511c62a06aa9f9c752f3ee61d73cf435ce4a1481f46f8348b21b82ec5a571a36ce4022dbf86b80843b9aca0082520894bdfeff9a1f4a1bdf483d680046344316019c58cf880de0a39a35d9b000801ba08327c70c73d0fcad956f760204e41c026c623bea1e38c1ca00930bd63d0a2384a068c4aba127f09d765dec74299de6a741e89ab792f630e6f827ea17f43f055400c0")
	var block ethtypes.Block
	_ = rlp.DecodeBytes(blockEnc, &block)

	listener := &dynamic.Listener{ExternalTxEnabled: false}

	c := &dynamic.Chain{
		UUID:     "test-c",
		URL:      "test-url",
		ChainID:  big.NewInt(1),
		Listener: listener,
	}

	// Test 1: envelope store is empty so producer should not be called thus no error
	receipt := &ethtypes.Receipt{
		TxHash:          ethcommon.HexToHash("0xdead"),
		ContractAddress: ethcommon.HexToAddress("0xbabde"),
	}
	err := hk.AfterNewBlock(context.Background(), c, &block, []*ethtypes.Receipt{receipt})
	assert.NoError(t, err, "#1 AfterNewBlock should not error")

	// Test 2: we store envelope on envelope store
	_, _ = store.Store(context.Background(), &evlpstore.StoreRequest{
		Envelope: tx.NewBuilder().SetID(uuid.NewV4().String()).SetChainID(c.ChainID).SetTxHash(receipt.TxHash).TxEnvelopeAsRequest(),
	})
	producer.ExpectSendMessageAndSucceed()
	err = hk.AfterNewBlock(context.Background(), c, &block, []*ethtypes.Receipt{receipt})
	assert.NoError(t, err, "#2 AfterNewBlock should not error")

	// Test 3: producer fails
	producer.ExpectSendMessageAndFail(fmt.Errorf("test-producer-error"))
	err = hk.AfterNewBlock(context.Background(), c, &block, []*ethtypes.Receipt{receipt})
	assert.Error(t, err, "#3 AfterNewBlock should error")
	assert.Equal(t, "test-producer-error", err.Error(), "#3 AfterNewBlock error message should be correct")
}
