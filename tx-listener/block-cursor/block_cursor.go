package blockcursor

import (
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/types"
)

// Cursor is an interface for a cursor object reading from a chain
type Cursor interface {
	// ChainID returns the chain ID the cursor is applied on
	ChainID() *big.Int

	// Current returns element the cursor is pointing on
	Blocks() <-chan *types.TxListenerBlock

	// Err returns a possible error met by the cursor when calling Next
	Errors() <-chan *types.TxListenerError

	// Close cursor
	Close()
}
