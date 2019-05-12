package listener

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/handler"
)

// TxListener is a transaction listener that allows to listen to multiple blockchains simultaneously
type TxListener interface {
	// Listen start listening for chains registered on a an Ethereum Client

	// It starts a blocking ListenerSession through the ListenerHandler.
	//
	// The life-cycle of a session is represented by the following steps:
	//
	// 1. The listener retrieve all Networks registered on the Ethereum client.
	// 2. Before processing starts, the handler's Setup() hook is called to notify the user
	//    of the Chains and allow any necessary preparation or alteration of state.
	// 3. For each of the Chains the handler's Listen() function is then called
	//    in a separate goroutine which requires it to be thread-safe. Any state must be carefully protected
	//    from concurrent reads/writes.
	// 4. The session will persist until one of the Listen() functions exits when the parent context is canceled
	// 5. Once all the Listen() loops have exited, the handler's Cleanup() hook is called
	//    to allow the user to perform any final tasks
	Listen(ctx context.Context, chains []*big.Int, h handler.TxListenerHandler) error

	// Close release listener resources
	Close()
}
