package handlers

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

// EthCrediter is an interface for crediting an account with ether
type EthCrediter interface {
	Credit(chainID *big.Int, a common.Address, value *big.Int) error
}

// CtxToCreditMessage create a Samara message corresponding to Eth Credit
type CtxToCreditMessage func(chainID *big.Int, a common.Address, value *big.Int) *sarama.ProducerMessage

// SaramaCrediter allow to credit an account by re-entering a message in core-stack
type SaramaCrediter struct {
	p sarama.SyncProducer

	makeCreditMsg CtxToCreditMessage
}

// Credit send a Credit message to a Kafka queue
func (c *SaramaCrediter) Credit(chainID *big.Int, a common.Address, value *big.Int) error {
	msg := c.makeCreditMsg(chainID, a, value)

	_, _, err := c.p.SendMessage(msg)
	return err
}

// EthCreditController is an interface to control if a credit should append
type EthCreditController interface {
	ShouldCredit(chainID *big.Int, a common.Address, value *big.Int) (*big.Int, bool)
}

// SimpleCreditController applies basic controls
type SimpleCreditController struct {
	cfg *SimpleCreditControllerConfig

	mux            *StripeMutex
	lastAuthorized *sync.Map
}

// NewSimpleCreditController creates a SimpleCreditController
func NewSimpleCreditController(cfg *SimpleCreditControllerConfig, stripes uint) *SimpleCreditController {
	return &SimpleCreditController{
		cfg: cfg,
		mux: NewStripeMutex(stripes),
		lastAuthorized: &sync.Map{},
	}
}

// ShouldCredit determines if a creadit should append
func (c *SimpleCreditController) ShouldCredit(chainID *big.Int, a common.Address, value *big.Int) (*big.Int, bool) {
	_, ok := c.cfg.blackList[a.Hex()]
	if ok {
		// Do not credit if address is in black list
		return nil, false
	}

	currentBalance, err := c.cfg.balanceAt(chainID, a)
	if err != nil {
		// Do not credit if we are not able to interogate client
		return nil, false
	}

	if currentBalance.Cmp(c.cfg.maxBalance) >= 0 {
		// Do not credit if current balance is higher than maximum authorized balance
		return nil, false
	}

	key := computeKey(chainID, a)
	c.mux.Lock(key)
	defer c.mux.Unlock(key)
	lastAuthorized, _ := c.lastAuthorized.LoadOrStore(key, time.Time{})
	if time.Now().Sub(lastAuthorized.(time.Time)) < c.cfg.creditDelay {
		// Do not credit if delay has not been respected
		return nil, false
	}

	// Update last authorization time
	c.lastAuthorized.Store(key, time.Now())

	return c.cfg.creditAmount, true
}

// BalanceAtFunc is a type for a function expected to return a balance
type BalanceAtFunc func(chainID *big.Int, a common.Address) (*big.Int, error)

// EthBalanceAt returns using an Ethereum Client
func EthBalanceAt(ec *infra.EthClient) BalanceAtFunc {
	return func(chainID *big.Int, a common.Address) (*big.Int, error) {
		return ec.BalanceAt(context.Background(), a, nil)
	}
}

// SimpleCreditControllerConfig is a config
type SimpleCreditControllerConfig struct {
	balanceAt    BalanceAtFunc
	creditAmount *big.Int
	maxBalance   *big.Int
	creditDelay  time.Duration
	blackList    map[string]struct{}
}

func txCost(tx *types.Tx) *big.Int {
	// Compute ETH cost of the transaction (based on formula cost = gasLimit*gasPrice)
	gas, cost := big.NewInt(0), big.NewInt(0)
	return cost.Mul(tx.GasPrice(), gas.SetUint64(tx.GasLimit()))
}

// Faucet creates a Faucet handler
func Faucet(crediter EthCrediter, controller EthCreditController) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		// Compute transaction cost
		cost := txCost(ctx.T.Tx())

		// Interogate credit controller
		amount, ok := controller.ShouldCredit(ctx.T.Chain().ID, *ctx.T.Sender().Address, cost)
		if !ok {
			// Credit invalid
			return
		}

		// We can credit
		crediter.Credit(ctx.T.Chain().ID, *ctx.T.Sender().Address, amount)
	}
}
