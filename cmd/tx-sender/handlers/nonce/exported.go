package nonce

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/nonce"
)

const component = "handler.nonce"

var (
	checker         engine.HandlerFunc
	recStatusSetter engine.HandlerFunc
	initOnce        = &sync.Once{}
)

// Init initialize Nonce Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		// Initialize the nonce manager
		nonce.Init(ctx)

		// Initialize the eth client
		ethclient.Init(ctx)

		if checker == nil {
			checker = Checker(nonce.GlobalManager(), ethclient.GlobalClient())
		}

		if recStatusSetter == nil {
			recStatusSetter = RecoveryStatusSetter(nonce.GlobalManager())
		}

		log.Infof("nonce: handlers checker & recovery status setter ready")
	})
}

// SetGlobalChecker sets global nonce checker
func SetGlobalChecker(h engine.HandlerFunc) {
	checker = h
}

// GlobalChecker returns global nonce checker handler
func GlobalChecker() engine.HandlerFunc {
	return checker
}

// SetGlobalRecoveryStatusSetter sets global nonce recovery status setter
func SetGlobalRecoveryStatusSetter(h engine.HandlerFunc) {
	recStatusSetter = h
}

// GlobalRecoveryStatusSetter returns global nonce recovery status setter
func GlobalRecoveryStatusSetter() engine.HandlerFunc {
	return recStatusSetter
}
