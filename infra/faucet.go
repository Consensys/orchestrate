package infra

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/striped-mutex.git"
)

// EthCrediter is an interface for crediting an account with ether
type EthCrediter interface {
	Credit(chainID *big.Int, a common.Address, value *big.Int) error
}

// MakeCreditMessageFunc is function expected to create Samara message corresponding to Eth Credit
type MakeCreditMessageFunc func(chainID *big.Int, a common.Address, value *big.Int) *sarama.ProducerMessage

// SaramaCrediter allow to credit an account by re-entering a message in core-stack
type SaramaCrediter struct {
	p sarama.SyncProducer

	makeCreditMsg MakeCreditMessageFunc
}

// NewSaramaCrediter creates a SaramaCrediter
func NewSaramaCrediter(p sarama.SyncProducer, makeCreditMsg MakeCreditMessageFunc) *SaramaCrediter {
	return &SaramaCrediter{p, makeCreditMsg}
}

// Credit send a Credit message to a Kafka queue
func (c *SaramaCrediter) Credit(chainID *big.Int, a common.Address, value *big.Int) error {
	msg := c.makeCreditMsg(chainID, a, value)
	_, _, err := c.p.SendMessage(msg)
	return err
}

// SimpleCreditController applies basic controls
type SimpleCreditController struct {
	cfg *SimpleCreditControllerConfig

	mux            *stripedmutex.StripedMutex
	lastAuthorized *sync.Map
}

// NewSimpleCreditController creates a SimpleCreditController
func NewSimpleCreditController(cfg *SimpleCreditControllerConfig, stripes uint) *SimpleCreditController {
	return &SimpleCreditController{
		cfg:            cfg,
		mux:            stripedmutex.New(stripes),
		lastAuthorized: &sync.Map{},
	}
}

// ShouldCredit determines if a creadit should append
func (c *SimpleCreditController) ShouldCredit(chainID *big.Int, a common.Address, value *big.Int) (*big.Int, bool) {
	_, ok := c.cfg.BlackList[a.Hex()]
	if ok {
		// Do not credit if address is in black list
		return nil, false
	}

	currentBalance, err := c.cfg.BalanceAt(chainID, a)
	if err != nil {
		// Do not credit if we are not able to interogate client
		return nil, false
	}

	if currentBalance.Cmp(c.cfg.MaxBalance) >= 0 {
		// Do not credit if current balance is higher than maximum authorized balance
		return nil, false
	}

	key := computeKey(chainID, &a)
	c.mux.Lock(key)
	defer c.mux.Unlock(key)
	lastAuthorized, _ := c.lastAuthorized.LoadOrStore(key, time.Time{})
	if time.Now().Sub(lastAuthorized.(time.Time)) < c.cfg.CreditDelay {
		// Do not credit if delay has not been respected
		return nil, false
	}

	// Update last authorization time
	c.lastAuthorized.Store(key, time.Now())

	return c.cfg.CreditAmount, true
}

// BalanceAtFunc is a type for a function expected to return a balance
type BalanceAtFunc func(chainID *big.Int, a common.Address) (*big.Int, error)

// NewEthBalanceAt returns using an Ethereum Client
func NewEthBalanceAt(ec *EthClient) BalanceAtFunc {
	return func(chainID *big.Int, a common.Address) (*big.Int, error) {
		return ec.BalanceAt(context.Background(), a, nil)
	}
}

// SimpleCreditControllerConfig is a config
type SimpleCreditControllerConfig struct {
	BalanceAt    BalanceAtFunc
	CreditAmount *big.Int
	MaxBalance   *big.Int
	CreditDelay  time.Duration
	BlackList    map[string]struct{}
}
