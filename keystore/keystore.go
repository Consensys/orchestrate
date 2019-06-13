package keystore

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
)

// KeyStore is an interface implemented by module that are able to perform signature transactions
type KeyStore interface {
	// SignTx signs a transaction
	SignTx(chain *chain.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, txHash *ethcommon.Hash, err error)

	// SignMsg sign a message TODO: what is the EIP?
	SignMsg(a ethcommon.Address, msg string) (rsv []byte, hash *ethcommon.Hash, err error) //TODO: do not forget to add prefix

	// SignRawHash sign a bytes
	SignRawHash(a ethcommon.Address, hash []byte) (rsv []byte, err error)

	// GenerateWallet creates a wallet
	GenerateWallet() (add *ethcommon.Address, err error)

	// ImportPrivateKey creates a wallet
	ImportPrivateKey(priv string) (err error)
}
