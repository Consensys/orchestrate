package listener

import (
	"context"
	"math/big"
	"sync"
	"time"

	backoff "github.com/cenkalti/backoff"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// EthClient is a minimal Ethereum Client interface required by a TxListener
type EthClient interface {
	// BlockByNumber retrieve a block by its number
	BlockByNumber(ctx context.Context, chainID *big.Int, number *big.Int) (*types.Block, error)

	// TransactionReceipt retrieve a transaction receipt using its hash
	TransactionReceipt(ctx context.Context, chainID *big.Int, txHash common.Hash) (*types.Receipt, error)

	// HeaderByNumber retrieves a block header
	HeaderByNumber(ctx context.Context, chainID *big.Int, number *big.Int) (*types.Header, error)
}

// TxListenerEthClient is an Ethereum with retry facilities
type TxListenerEthClient struct {
	c EthClient

	pool *sync.Pool
}

func newBackOff() backoff.BackOff {
	return &backoff.ExponentialBackOff{
		InitialInterval:     viper.GetDuration("ethclient.retry.initinterval"),
		RandomizationFactor: viper.GetFloat64("ethclient.retry.randomfactor"),
		Multiplier:          viper.GetFloat64("ethclient.retry.multiplier"),
		MaxInterval:         viper.GetDuration("ethclient.retry.maxinterval"),
		MaxElapsedTime:      viper.GetDuration("ethclient.retry.maxelapsedtime"),
		Clock:               backoff.SystemClock,
	}
}

// NewEthClient creates an Ethereum client compatible with a TxListener
func NewEthClient(ec EthClient) EthClient {
	return &TxListenerEthClient{
		c: ec,
		pool: &sync.Pool{
			New: func() interface{} { return newBackOff() },
		},
	}
}

// HeaderByNumber returns a block from the current canonical chain. If number is
// nil, the latest known header is returned.
func (ec *TxListenerEthClient) HeaderByNumber(ctx context.Context, chainID *big.Int, number *big.Int) (*types.Header, error) {
	var res *types.Header
	// Try retrieving header with backoff strategy
	bckoff := backoff.WithContext(ec.pool.Get().(backoff.BackOff), ctx)
	defer ec.pool.Put(bckoff)

	err := backoff.RetryNotify(
		func() error {
			header, err := ec.c.HeaderByNumber(ctx, chainID, number)
			if err != nil {
				return err
			}
			res = header
			return nil
		},
		bckoff,
		func(err error, duration time.Duration) {
			log.WithFields(log.Fields{
				"Chain":       chainID.Text(16),
				"BlockNumber": number.Text(10),
			}).Warnf("tx-listener: error retrieving header, retrying in %v...", duration)
		},
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// BlockByNumber returns a block from the current canonical chain. If number is
// nil, the latest known header is returned.
func (ec *TxListenerEthClient) BlockByNumber(ctx context.Context, chainID *big.Int, number *big.Int) (*types.Block, error) {
	var res *types.Block
	// Try retrieving block with backoff strategy
	bckoff := backoff.WithContext(ec.pool.Get().(backoff.BackOff), ctx)
	defer ec.pool.Put(bckoff)

	err := backoff.RetryNotify(
		func() error {
			block, err := ec.c.BlockByNumber(ctx, chainID, number)
			if err != nil {
				return err
			}
			res = block
			return nil
		},
		bckoff,
		func(err error, duration time.Duration) {
			log.WithError(err).WithFields(log.Fields{
				"Chain":       chainID.Text(16),
				"BlockNumber": number.Text(10),
			}).Warnf("tx-listener: error retrieving block, retrying in %v...", duration)
		},
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (ec *TxListenerEthClient) TransactionReceipt(ctx context.Context, chainID *big.Int, txHash common.Hash) (*types.Receipt, error) {
	var res *types.Receipt

	// Try retrieving receipt with backoff strategy
	bckoff := backoff.WithContext(ec.pool.Get().(backoff.BackOff), ctx)
	defer ec.pool.Put(bckoff)

	err := backoff.RetryNotify(
		func() error {
			receipt, err := ec.c.TransactionReceipt(ctx, chainID, txHash)
			if err != nil {
				return err
			}
			res = receipt
			return nil
		},
		bckoff,
		func(err error, duration time.Duration) {
			log.WithError(err).WithFields(log.Fields{
				"Chain":  chainID.Text(16),
				"TxHash": txHash.Hex(),
			}).Warnf("tx-listener: error retrieving receipt, retrying in %v...", duration)
		},
	)

	if err != nil {
		return nil, err
	}

	ec.pool.Put(bckoff)

	return res, nil
}
