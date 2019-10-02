package handler

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/types"
)

// ListenerSession
type TxListenerSession interface {
	Chains() []*big.Int
}

// ChainListener is a listener on a given chain
type ChainListener interface {
	// ChainID returns ID of the chain being listen
	ChainID() *big.Int

	// InitialPosition returns position from which the listener started
	InitialPosition() (int64, int64)

	// Receipts returns a channel of ordered receipts
	Receipts() <-chan *types.TxListenerReceipt

	// Blocks returns a channel of Blocks are they are mined
	Blocks() <-chan *types.TxListenerBlock

	// Errors returns a channel of Errors encountered while listening
	Errors() <-chan *types.TxListenerError

	// Context associated to the current
	Context() context.Context
}

// TxListenerHandler instances are used to handle individual .
// It also provides hooks to allow you to trigger custom logic before or after the listening loop(s).
//
// PLEASE NOTE that handlers are likely to be called from several goroutines concurrently,
// ensure that all state is safely protected against race conditions.
type TxListenerHandler interface {
	// Setup is run at the beginning of a new session, before Listen.
	Setup(TxListenerSession) error

	// GetInitialPosition is called right before listen to get initial position to start listening
	// from on a given chain
	// To start from latest block set blockNumber -1
	GetInitialPosition(chainID *big.Int) (blockNumber int64, txIndex int64, err error)

	// Cleanup is run at the end of a session, once all Listen goroutines have exited
	Cleanup(TxListenerSession) error

	// Listen must start a consumer loop of ChainListener Receipts() (and optionally Blocks() and Errors())
	// Once the Receipts() channel is closed, the Handler must finish its processing
	// loop and exit.
	Listen(TxListenerSession, ChainListener) error
}
