package services

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// KeyStore is an interface implemented by module that are able to perform signature transactions
type KeyStore interface {
	SignTx(chain *types.Chain, a common.Address, tx *ethtypes.Transaction) (raw []byte, txHash *common.Hash, err error)
	SignMsg(a common.Address, msg string) (rsv []byte, hash *common.Hash, err error) //TODO: do notforget to add prefix
	SignRawHash(a common.Address, hash []byte) (rsv []byte, err error)
	GenerateWallet() (add *common.Address, err error)
}
