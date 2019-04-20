package mock

import (
	"context"
	"sync"
)

var (
	fct      *Faucet
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		// Initialize Faucet
		fct = NewFaucet()
	})
}

// GlobalFaucet returns global Sarama Faucet
func GlobalFaucet() *Faucet {
	return fct
}

// SetGlobalFaucet sets global Sarama Faucet
func SetGlobalFaucet(faucet *Faucet) {
	initOnce.Do(func() {
		fct = faucet
	})
}
