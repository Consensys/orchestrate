package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/secretstore/services"
)

// Wallet is a container
type Wallet struct {
	address common.Address
	priv    *ecdsa.PrivateKey

	sec services.SecretStore
}

// NewWallet construct a wallet object
func NewWallet(sec services.SecretStore) *Wallet {
	return &Wallet{
		sec: sec,
	}
}

// Generate create a keypair and set the result in the wallet
func (w *Wallet) Generate() error {
	prv, err := crypto.GenerateKey()
	if err != nil {
		return errors.InternalError(err.Error()).SetComponent(component)
	}
	w.priv = prv
	pub := prv.PublicKey
	w.address = crypto.PubkeyToAddress(pub)
	return nil
}

// FromPrivateKey creates a new wallet from a given private key
func (w *Wallet) FromPrivateKey(priv string) error {
	prv, err := crypto.HexToECDSA(priv)
	if err != nil {
		return errors.InvalidFormatError("invalid hex private key %v...%v", priv[:5], priv[len(priv)-5:]).SetComponent(component)
	}
	w.priv = prv
	pub := w.priv.PublicKey
	w.address = crypto.PubkeyToAddress(pub)
	return nil
}

// Store saves wallet information to secret store
func (w *Wallet) Store() error {
	if err := w.sec.Store(
		w.address.Hex(),
		hex.EncodeToString(crypto.FromECDSA(w.priv)),
	); err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

// Load wallets values by fetching wallet secret store
func (w *Wallet) Load(a *common.Address) (err error) {
	w.address = *a
	priv, ok, err := w.sec.Load(a.Hex())
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

// Priv returns the private key of Wallet
func (w *Wallet) Priv() *ecdsa.PrivateKey {
	return w.priv
}

// Address returns the address of the wallet
func (w *Wallet) Address() *common.Address {
	return &w.address
}
