package noncechecker

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce"
)

const component = "handler.nonce.checker"

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

		// Create recovery tracker
		tracker := NewRecoveryTracker()

		conf := NewConfig()
		if checker == nil {
			checker = Checker(conf, nonce.GlobalManager(), ethclient.GlobalClient(), tracker)
		}

		if recStatusSetter == nil {
			recStatusSetter = RecoveryStatusSetter(nonce.GlobalManager(), tracker)
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
