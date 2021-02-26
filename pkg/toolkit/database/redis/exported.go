package redis

import (
	"sync"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	healthz "github.com/heptiolabs/healthcheck"
	"github.com/spf13/viper"
)

const component = "database.redis"

var (
	nm       *Client
	initOnce = &sync.Once{}
	checker  healthz.Check
)

// Init initializes Nonce
func Init() {
	initOnce.Do(func() {
		if nm != nil {
			return
		}

		cfg := NewConfig(viper.GetViper())
		logger := log.NewLogger().SetComponent(component).WithField("host", cfg.URL())

		pool, err := NewPool(cfg)
		if err != nil {
			logger.Fatalf("could not connect to server")
		}

		// Initialize Nonce
		nm = NewClient(pool, cfg)

		checker = healthz.Timeout(nm.Ping, time.Second*2)
		logger.Info("ready")
	})
}

// GlobalClient returns global NonceManager
func GlobalClient() *Client {
	return nm
}

func GlobalChecker() healthz.Check {
	return checker
}

// SetGlobalNonce sets global NonceManager
func SetGlobalClient(m *Client) {
	nm = m
}
