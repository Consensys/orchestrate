package account

import (
	"crypto/ecdsa"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

const component = "keystore.account"

type Account interface {
	Generate() error
	FromPrivateKey(string) error
	Priv() *ecdsa.PrivateKey
	Address() ethcommon.Address
	SetPriv(*ecdsa.PrivateKey)
	SetAddress(ethcommon.Address)
}
