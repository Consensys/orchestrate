package types

import (
	"fmt"
	"math/big"

	"github.com/ConsenSys/orchestrate/pkg/engine"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// TxListenerReceipt contains useful information about a receipt
type TxListenerReceipt struct {
	// Chain receipt has been read from
	ChainID *big.Int

	// Go-Ethereum receipt
	ethtypes.Receipt

	// Position of the receipt
	BlockHash   ethcommon.Hash
	BlockNumber int64
	TxHash      ethcommon.Hash
	TxIndex     uint64
}

// make TxListenerReceipt match engine.Msg interface
func (m *TxListenerReceipt) Entrypoint() string    { return "" }
func (m *TxListenerReceipt) Header() engine.Header { return &header{} }
func (m *TxListenerReceipt) Value() []byte         { return []byte{} }
func (m *TxListenerReceipt) Key() []byte           { return []byte{} }

type header struct{}

func (h *header) Add(key, value string) {}
func (h *header) Del(key string)        {}
func (h *header) Get(key string) string { return "" }
func (h *header) Set(key, value string) {}

// TxListenerBlock contains data about a block
type TxListenerBlock struct {
	// Chain block has been read from
	ChainID *big.Int

	// Go-Ethereum block
	ethtypes.Block

	// Ordered receipts for every transaction in the block
	Receipts []*TxListenerReceipt
}

// Copy creates a deep copy of a block to prevent side effects
func (b *TxListenerBlock) Copy() *TxListenerBlock {
	return &TxListenerBlock{
		ChainID:  big.NewInt(0).Set(b.ChainID),
		Block:    *b.WithBody(b.Transactions(), b.Uncles()),
		Receipts: make([]*TxListenerReceipt, len(b.Receipts)),
	}
}

// TxListenerError is what is provided to the user when an error occurs.
// It wraps an error and includes the chain UUID
type TxListenerError struct {
	// Network UUID the error occurred on
	ChainID *big.Int

	// Error
	Err error
}

func (e TxListenerError) Error() string {
	return fmt.Sprintf("tx-listener: error while listening on chain %s: %s", hexutil.EncodeBig(e.ChainID), e.Err)
}

// TxListenerErrors is a type that wraps a batch of errors and implements the Error interface.
type TxListenerErrors []*TxListenerErrors

func (e TxListenerErrors) Error() string {
	return fmt.Sprintf("tx-listener: %d errors while while listening", len(e))
}

// Progress holds information about listener progress
type Progress struct {
	CurrentBlock int64 // Current block number where the listener is
	TxIndex      int64 // Current txIndex where the listener is
	HighestBlock int64 // Highest alleged block number in the chain
}
