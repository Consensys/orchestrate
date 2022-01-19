package account

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// TODO: Remove this Account entity and use the one in pkg/types

type Account struct {
	Address common.Address
	priv    *ecdsa.PrivateKey
}

// NewAccount create a keypair and set the result in the account
func NewAccount() (*Account, error) {
	prv, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	return &Account{
		priv:    prv,
		Address: crypto.PubkeyToAddress(prv.PublicKey),
	}, nil
}

// NewAccountFromPrivateKey creates a new account from a given private key
func NewAccountFromPrivateKey(priv string) (*Account, error) {
	prv, err := crypto.HexToECDSA(priv[2:])
	if err != nil {
		return nil, err
	}

	return &Account{
		priv:    prv,
		Address: crypto.PubkeyToAddress(prv.PublicKey),
	}, nil
}

func (a *Account) Priv() []byte {
	return a.priv.D.Bytes()
}
