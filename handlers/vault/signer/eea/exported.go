package eea

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/onetimekey"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore"
)

const component = "handler.signer.eea"

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Gas Estimator Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		utils.InParallel(
			// Initialize keystore
			func() { keystore.Init(ctx) },
			// Initialize OneTimeKey Signer
			func() { onetimekey.Init(ctx) },
			// Initialize OneTimeKey Signer
			func() { ethclient.Init(ctx) },
		)

		// Create Handler
		handler = Signer(keystore.GlobalKeyStore(), onetimekey.GlobalKeyStore(), ethclient.GlobalClient())

		log.Infof("eea signer: handler ready")
	})
}

// SetGlobalHandler sets global Gas Estimator Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalKeyStore returns global Gas Estimator handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
