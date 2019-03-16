package keystore

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

// KeyStore is an interface implemented by module that are able to perform signature transactions
type KeyStore interface {
	SignTx(chain *common.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, txHash *ethcommon.Hash, err error)
	SignMsg(a ethcommon.Address, msg string) (rsv []byte, hash *ethcommon.Hash, err error) //TODO: do notforget to add prefix
	SignRawHash(a ethcommon.Address, hash []byte) (rsv []byte, err error)
	GenerateWallet() (add *ethcommon.Address, err error)
}
