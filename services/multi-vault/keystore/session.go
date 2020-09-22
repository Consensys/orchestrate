package keystore

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/session"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/services"
)

type accountManager struct {
	SecretStore services.SecretStore
}

// NewKeyStore construct a BaseKeyStore from a client
func NewKeyStore(secretStore services.SecretStore) keystore.KeyStore {
	return keystore.NewKeyStore(newAccountManager(secretStore))
}

func newAccountManager(secretStore services.SecretStore) session.AccountManager {
	return &accountManager{
		SecretStore: secretStore,
	}
}

func (am *accountManager) SigningSession(ctx context.Context, addr ethcommon.Address) (session.SigningSession, error) {
	// Creates a new signing session
	sess := session.NewSigningSession()
	w := NewAccount(am.SecretStore)
	err := w.Load(ctx, addr)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	log.WithContext(ctx).Debugf("new multi-vault signing session %s", addr.String())
	err = sess.SetAccount(w)
	return sess, err
}

// GenerateAccount create and stores a new account in the vault
func (am *accountManager) GenerateAccount(ctx context.Context) (ethcommon.Address, error) {
	w := NewAccount(am.SecretStore)
	err := w.Generate()
	if err != nil {
		return ethcommon.Address{}, errors.FromError(err).ExtendComponent(component)
	}

	err = w.Store(ctx)
	if err != nil {
		return ethcommon.Address{}, errors.FromError(err).ExtendComponent(component)
	}

	log.WithContext(ctx).WithField("address", w.Address().String()).
		WithField("tenant", multitenancy.TenantIDFromContext(ctx)).
		Debug("new multi-vault account generated")
	return w.Address(), nil
}

// ImportPrivateKey adds a private key in the vault
// TODO: this is Unsafe and should be removed soon
func (am *accountManager) ImportPrivateKey(ctx context.Context, priv string) error {
	w := NewAccount(am.SecretStore)
	err := w.FromPrivateKey(priv)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	err = w.Store(ctx)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	log.WithContext(ctx).Debugf("new multi-vault account imported %s", w.Address().String())

	return nil
}
