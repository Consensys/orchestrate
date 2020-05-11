package onetimekey

import (
	"context"
	"strings"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/account"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/session"
)

type accountManager struct {
	accounts map[string]account.Account
	mux      *sync.RWMutex
}

// NewKeyStore construct a BaseKeyStore from a client
func NewKeyStore() keystore.KeyStore {
	return keystore.NewKeyStore(newAccountManager())
}

func newAccountManager() session.AccountManager {
	return &accountManager{
		accounts: make(map[string]account.Account),
		mux:      &sync.RWMutex{},
	}
}

func (am *accountManager) SigningSession(ctx context.Context, addr ethcommon.Address) (session.SigningSession, error) {
	acc := am.getAccount(addr)
	if acc == nil {
		log.WithContext(ctx).Errorf("account %s was not found", addr.String())
		return nil, errors.DataError("one-time-key account %s was not found", addr.String())
	}

	sess := session.NewSigningSession()
	err := sess.SetAccount(acc)
	log.WithContext(ctx).Debugf("created new one-time-key signing session %s", addr.String())
	return sess, err

}

// GenerateAccount create and stores a new account in the vault
func (am *accountManager) GenerateAccount(ctx context.Context) (ethcommon.Address, error) {
	w := account.NewAccount()
	err := w.Generate()
	if err != nil {
		return ethcommon.Address{}, errors.FromError(err).ExtendComponent(component)
	}

	log.WithContext(ctx).Debugf("new one-time-key generated %s", w.Address().String())
	err = am.setAccount(w)
	if err != nil {
		return ethcommon.Address{}, errors.FromError(err).ExtendComponent(component)
	}

	return w.Address(), nil
}

// ImportPrivateKey adds a private key in the vault
// TODO: this is Unsafe and should be removed soon
func (am *accountManager) ImportPrivateKey(ctx context.Context, priv string) error {
	w := account.NewAccount()
	err := w.FromPrivateKey(priv)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	err = am.setAccount(w)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	log.WithContext(ctx).Debugf("new one-time-key imported %s", w.Address().String())
	return nil
}

func (am *accountManager) setAccount(w account.Account) error {
	am.mux.Lock()
	defer am.mux.Unlock()
	arrKey := strings.ToLower(w.Address().String())
	if _, ok := am.accounts[arrKey]; ok {
		return errors.AlreadyExistsError("account %s was already used", arrKey)
	}

	am.accounts[arrKey] = w
	return nil
}

func (am *accountManager) getAccount(addr ethcommon.Address) account.Account {
	am.mux.RLock()
	defer am.mux.RUnlock()
	arrKey := strings.ToLower(addr.String())
	if acc, ok := am.accounts[arrKey]; ok {
		return acc
	}

	return nil
}
