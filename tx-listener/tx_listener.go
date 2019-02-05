package listener

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TxListenerEthClient is a minimal EthClient interface required by a TxListener
type TxListenerEthClient interface {
	// BlockByNumber retrieve a block by its number
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)

	// TransactionReceipt retrieve a transaction receipt given a transaction hash
	TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error)
}

// TxListener is an interface to listen to a Ethereum blockchain activity
type TxListener interface {
	// Blocks return a read channel of blocks
	Blocks() <-chan *types.Block

	// Receipts return a read channel of transaction receipts
	Receipts() <-chan *types.Receipt

	// Errors return a read channel of errors
	Errors() <-chan error

	// Set allows to set TxListener position
	Set(block uint64, txIndex uint64)

	// Close stops listener
	Close()
}

type receiptResponse struct {
	txHash common.Hash
	err    chan error
	res    chan *types.Receipt
}

type txListener struct {
	ec  TxListenerEthClient
	cfg *Config

	receipts chan *types.Receipt
	errors   chan error
	blocks   chan *types.Block
	pos      uint64

	bl        BlockListener
	responses chan *receiptResponse

	closeOnce *sync.Once
	closed    chan struct{}
}

// NewTxListener creates a new transaction listener
func NewTxListener(ec TxListenerEthClient, conf *Config) TxListener {
	// Instantiate txListener
	cur := newBlockCursor(ec)
	bl := newBlockListener(cur, conf)
	l := &txListener{
		ec:        ec,
		cfg:       conf,
		receipts:  make(chan *types.Receipt),
		errors:    make(chan error),
		blocks:    make(chan *types.Block),
		bl:        bl,
		responses: make(chan *receiptResponse, conf.TxListener.MaxReceiptCount),
		closeOnce: &sync.Once{},
		closed:    make(chan struct{}),
	}

	// Start feeding txListener
	go l.feeder()
	go l.blockDispatcher()
	go l.errorDispatcher()

	// Start feeding block consumer
	go bl.feeder()

	return l
}

// Warning: you should not call SetPos while consuming listener channels
func (l *txListener) Set(block uint64, txIndex uint64) {
	l.bl.Set(big.NewInt(int64(block)))
	l.pos = txIndex
}

func (l *txListener) Blocks() <-chan *types.Block {
	return l.blocks
}

func (l *txListener) Receipts() <-chan *types.Receipt {
	return l.receipts
}

func (l *txListener) Errors() <-chan error {
	return l.errors
}

func (l *txListener) errorDispatcher() {
	for err := range l.bl.Errors() {
		if l.cfg.TxListener.Return.Errors {
			l.errors <- err
		}
	}
	close(l.errors)
}

func (l *txListener) blockDispatcher() {
dispatchLoop:
	for block := range l.bl.Blocks() {
		// A new block has been mined
		// Send block to blocks channel
		if l.cfg.TxListener.Return.Blocks {
			l.blocks <- block
		}

		// Retrieve transaction receipt for every transactions
		for _, tx := range block.Transactions()[l.pos:] {
			select {
			case <-l.closed:
				break dispatchLoop
			default:
				// We retrieve receipt concurrently
				l.responses <- l.getReceipt(tx)
			}
		}
		// Reinitialize position
		l.pos = 0
	}
	close(l.responses)
	close(l.blocks)
	l.bl.Close()
}

func (l *txListener) feeder() {
	for response := range l.responses {
		// Retrieve receipt response and wait for it to complete
		select {
		case err := <-response.err:
			if l.cfg.TxListener.Return.Errors {
				l.errors <- err
			}
		case receipt := <-response.res:
			l.receipts <- receipt
		}
	}
	close(l.receipts)
}

func (l *txListener) getReceipt(tx *types.Transaction) *receiptResponse {
	response := &receiptResponse{
		txHash: tx.Hash(),
		err:    make(chan error, 1),
		res:    make(chan *types.Receipt, 1),
	}
	// Retrieve receipt in a parallel goroutine
	go func() {
		receipt, err := l.ec.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			response.err <- err
		} else {
			response.res <- receipt
		}
	}()
	return response
}

func (l *txListener) Close() {
	l.closeOnce.Do(func() {
		close(l.closed)
	})
}
