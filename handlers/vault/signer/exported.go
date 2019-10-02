package signer

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/vault/signer/eea"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/vault/signer/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/vault/signer/tessera"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
)

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

		common.InParallel(
			// Initialize EEA signer
			func() { eea.Init(ctx) },
			// Initialize Tessera Signer
			func() { tessera.Init(ctx) },
			// Initialize Public Ethereum Signer
			func() { ethereum.Init(ctx) },
		)

		// Create Handler
		handler = engine.CombineHandlers(
			TxSigner(eea.GlobalHandler(), ethereum.GlobalHandler(), tessera.GlobalHandler()),
		)

		log.Infof("signer: handler ready")
	})
}

// SetGlobalHandler sets global Gas Estimator Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Gas Estimator handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
