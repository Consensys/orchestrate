package account

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/services"
)

// Account is a container
type Account struct {
	address common.Address
	priv    *ecdsa.PrivateKey

	sec services.SecretStore
}

// NewAccount construct an account object
func NewAccount(sec services.SecretStore) *Account {
	return &Account{
		sec: sec,
	}
}

// Generate create a keypair and set the result in the account
func (w *Account) Generate() error {
	prv, err := crypto.GenerateKey()
	if err != nil {
		return errors.InternalError(err.Error()).SetComponent(component)
	}
	w.priv = prv
	pub := prv.PublicKey
	w.address = crypto.PubkeyToAddress(pub)
	return nil
}

// FromPrivateKey creates a new account from a given private key
func (w *Account) FromPrivateKey(priv string) error {
	prv, err := crypto.HexToECDSA(priv)
	if err != nil {
		return errors.InvalidFormatError("invalid hex private key %v...%v", priv[:5], priv[len(priv)-5:]).SetComponent(component)
	}
	w.priv = prv
	pub := w.priv.PublicKey
	w.address = crypto.PubkeyToAddress(pub)
	return nil
}

// Store saves account information to secret store
func (w *Account) Store(ctx context.Context) error {
	if err := w.sec.Store(
		ctx,
		w.address.Hex(),
		hex.EncodeToString(crypto.FromECDSA(w.priv)),
	); err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

// Load accounts values by fetching account secret store
func (w *Account) Load(ctx context.Context, a *common.Address) (err error) {
	w.address = *a
	priv, ok, err := w.sec.Load(ctx, a.Hex())
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if !ok {
		return errors.NotFoundError("no key for account %q", a.Hex())
	}

	w.priv, err = crypto.HexToECDSA(priv)
	if err != nil {
		return errors.InvalidFormatError("invalid hex private key %v...%v", priv[:5], priv[len(priv)-5:]).SetComponent(component)
	}

	return nil
}

// Priv returns the private key of Account
func (w *Account) Priv() *ecdsa.PrivateKey {
	return w.priv
}

// Address returns the address of the account
func (w *Account) Address() *common.Address {
	return &w.address
}
