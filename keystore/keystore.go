package keystore

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
)

// KeyStore is an interface implemented by module that are able to perform signature transactions
type KeyStore interface {
	// SignTx signs a transaction
	SignTx(chain *chain.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, txHash *ethcommon.Hash, err error)

	// SignPrivateEEATx signs a private transaction
	SignPrivateEEATx(chain *chain.Chain, a ethcommon.Address, tx *ethtypes.Transaction, privateArgs *types.PrivateArgs) (raw []byte, txHash *ethcommon.Hash, err error)

	// SignPrivateTesseraTx signs a private transaction for Tessera transactions manager
	// Before calling this method, "data" field in the transaction should be replaced with the result
	// of the "storeraw" API call
	SignPrivateTesseraTx(chain *chain.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, txHash *ethcommon.Hash, err error)

	// SignMsg sign a message TODO: what is the EIP?
	SignMsg(a ethcommon.Address, msg string) (rsv []byte, hash *ethcommon.Hash, err error) //TODO: do not forget to add prefix

	// SignRawHash sign a bytes
	SignRawHash(a ethcommon.Address, hash []byte) (rsv []byte, err error)

	// GenerateWallet creates a wallet
	GenerateWallet() (add *ethcommon.Address, err error)

	// ImportPrivateKey creates a wallet
	ImportPrivateKey(priv string) (err error)
}

// ImportPrivateKey create new Key Store
func ImportPrivateKey(k KeyStore, pkeys []string) error {
	// Pre-Import Pkeys
	for _, pkey := range pkeys {
		err := k.ImportPrivateKey(pkey)
		if err != nil {
			return err
		}
	}

	return nil
}
