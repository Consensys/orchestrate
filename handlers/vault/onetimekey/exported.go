package onetimekey

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore"
)

const component = "handler.signer.ethereum"

var (
	ks       keystore.KeyStore
	initOnce = &sync.Once{}
)

// Init initialize Gas Estimator Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if ks != nil {
			return
		}

		ks = NewKeyStore()

		log.Infof("public one-time-key signer: handler ready")
	})
}

// SetGlobalHandler sets global Gas Estimator Handler
func SetGlobalKeyStore(k keystore.KeyStore) {
	ks = k
}

// GlobalKeyStore returns global Gas Estimator handler
func GlobalKeyStore() keystore.KeyStore {
	return ks
}
