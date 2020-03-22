package key

import (
	"context"
	"sync"

	"github.com/spf13/viper"
)

var (
	checker  *Key
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if checker != nil {
			return
		}

		checker = New(viper.GetString(APIKeyViperKey))
	})
}

// GlobalChecker returns global Checker
func GlobalChecker() *Key {
	return checker
}

// SetGlobalChecker sets global Checker
func SetGlobalChecker(c *Key) {
	checker = c
}
