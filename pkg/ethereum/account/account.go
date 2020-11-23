package account

import (
	"crypto/ecdsa"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

// account is a container
type account struct {
	address ethcommon.Address
	priv    *ecdsa.PrivateKey
}

// NewAccount construct an account object
func NewAccount() Account {
	return &account{}
}

// Generate create a keypair and set the result in the account
func (w *account) Generate() error {
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
func (w *account) FromPrivateKey(priv string) error {
	prv, err := crypto.HexToECDSA(priv)
	if err != nil {
		return errors.InvalidFormatError("invalid hex private key %v...%v", priv[:5], priv[len(priv)-5:]).SetComponent(component)
	}
	w.priv = prv
	pub := w.priv.PublicKey
	w.address = crypto.PubkeyToAddress(pub)
	return nil
}

// Priv returns the private key of account
func (w *account) Priv() *ecdsa.PrivateKey {
	return w.priv
}

// Address returns the address of the account
func (w *account) Address() ethcommon.Address {
	return w.address
}

func (w *account) SetPriv(priv *ecdsa.PrivateKey) {
	w.priv = priv
}

func (w *account) SetAddress(addr ethcommon.Address) {
	w.address = addr
}
