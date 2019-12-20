package keystore

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/chain"
)

// KeyStore is an interface implemented by module that are able to perform signature transactions
type KeyStore interface {
	// SignTx signs a transaction
	SignTx(ctx context.Context, chain *chain.Chain, a ethcommon.Address, tx *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error)

	// SignPrivateEEATx signs a private transaction
	SignPrivateEEATx(ctx context.Context, chain *chain.Chain, a ethcommon.Address, tx *ethtypes.Transaction, privateArgs *types.PrivateArgs) ([]byte, *ethcommon.Hash, error)

	// SignPrivateTesseraTx signs a private transaction for Tessera transactions manager
	// Before calling this method, "data" field in the transaction should be replaced with the result
	// of the "storeraw" API call
	SignPrivateTesseraTx(ctx context.Context, chain *chain.Chain, a ethcommon.Address, tx *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error)

	// SignMsg sign a message
	SignMsg(ctx context.Context, a ethcommon.Address, msg string) ([]byte, *ethcommon.Hash, error)

	// SignRawHash sign a bytes
	SignRawHash(a ethcommon.Address, hash []byte) ([]byte, error)

	// GenerateWallet creates a wallet
	GenerateWallet(ctx context.Context) (*ethcommon.Address, error)

	// ImportPrivateKey creates a wallet
	ImportPrivateKey(ctx context.Context, priv string) error
}

// ImportPrivateKey create new Key Store
func ImportPrivateKey(k KeyStore, pkeys []string) error {
	// Pre-Import Pkeys
	for _, pkey := range pkeys {
		ctx, key, err := multitenancy.SplitTenant(pkey)
		if err != nil {
			return err
		}
		err = k.ImportPrivateKey(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}
