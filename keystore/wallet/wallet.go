package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore"
)

// Wallet is a container
type Wallet struct {
	address common.Address
	priv    *ecdsa.PrivateKey

	sec secretstore.SecretStore
}

// NewWallet construct a wallet object
func NewWallet(sec secretstore.SecretStore) *Wallet {
	return &Wallet{
		sec: sec,
	}
}

// Generate create a keypair and set the result in the wallet
func (w *Wallet) Generate() error {
	prv, err := crypto.GenerateKey()
	if err != nil {
		return err
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
		return err
	}
	w.priv = prv
	pub := w.priv.PublicKey
	w.address = crypto.PubkeyToAddress(pub)
	return nil
}

// Store saves wallet information to secret store
func (w *Wallet) Store() (err error) {
	pathStr := fmt.Sprintf("%v/%v", "private_keys", w.address.Hex())
	return w.sec.Store(
		pathStr,
		hex.EncodeToString(crypto.FromECDSA(w.priv)),
	)
}

// Load wallets values by fetching wallet secret store
func (w *Wallet) Load(a *common.Address) (err error) {
	w.address = *a

	pathStr := fmt.Sprintf("%v/%v", "private_keys", a.Hex())
	priv, ok, err := w.sec.Load(pathStr)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("No pkey stored for account %q", a.Hex())
	}

	w.priv, err = crypto.HexToECDSA(priv)
	if err != nil {
		return fmt.Errorf("Could not deserialize %v...%v", priv[:5], priv[len(priv)-5:])
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
