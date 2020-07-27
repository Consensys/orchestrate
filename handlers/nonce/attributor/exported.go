package nonceattributor

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce"
)

const component = "handler.nonce.attributor"

var (
	handler    engine.HandlerFunc
	eeaHandler engine.HandlerFunc
	initOnce   = &sync.Once{}
)

// Init initialize Nonce Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize the nonce manager
		nonce.Init(ctx)

		// Initialize the eth client
		ethclient.Init(ctx)

		// Create Nonce handler
		handler = Nonce(nonce.GlobalManager(), ethclient.GlobalClient())

		// @TODO Remove once Orion tx is spit into two jobs
		//  https://app.zenhub.com/workspaces/orchestrate-5ea70772b186e10067f57842/issues/pegasyseng/orchestrate/253
		eeaHandler = EEANonce(nonce.GlobalManager(), ethclient.GlobalClient())

		log.Infof("%s: handler ready", component)
	})
}

// SetGlobalHandler sets global Faucet Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Faucet handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}

// GlobalHandler returns global Faucet handler
func GlobalEEAHandler() engine.HandlerFunc {
	return eeaHandler
}
