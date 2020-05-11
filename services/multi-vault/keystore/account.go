package keystore

import (
	"context"
	"encoding/hex"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/account"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/services"
)

// Account is a container
type Account struct {
	account.Account
	sec services.SecretStore
}

// NewAccount construct an account object
func NewAccount(sec services.SecretStore) *Account {
	return &Account{
		Account: account.NewAccount(),
		sec:     sec,
	}
}

// Store saves account information to secret store
func (w *Account) Store(ctx context.Context) error {
	if err := w.sec.Store(
		ctx,
		w.Address().Hex(),
		hex.EncodeToString(crypto.FromECDSA(w.Priv())),
	); err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

// Load accounts values by fetching account secret store
func (w *Account) Load(ctx context.Context, a ethcommon.Address) (err error) {
	w.SetAddress(a)
	priv, ok, err := w.sec.Load(ctx, a.Hex())
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if !ok {
		return errors.NotFoundError("no key for account %q", a.Hex())
	}

	npriv, err := crypto.HexToECDSA(priv)
	if err != nil {
		return errors.InvalidFormatError("invalid hex private key %v...%v", priv[:5], priv[len(priv)-5:]).SetComponent(component)
	}
	w.SetPriv(npriv)

	return nil
}
