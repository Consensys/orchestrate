package keystore

import (
	"context"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
)

//go:generate mockgen -source=exported.go -destination=mock/keystore.go -package=mock

const component = "keystore"

// KeyStore is an interface implemented by module that are able to perform signature transactions
type KeyStore interface {
	// SignTx signs a transaction
	SignTx(context.Context, *big.Int, ethcommon.Address, *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error)

	// SignPrivateEEATx signs a private transaction
	SignPrivateEEATx(context.Context, *big.Int, ethcommon.Address, *ethtypes.Transaction, *types.PrivateArgs) ([]byte, *ethcommon.Hash, error)

	// SignPrivateTesseraTx signs a private transaction for Tessera transactions manager
	// Before calling this method, "data" field in the transaction should be replaced with the result
	// of the "storeraw" API call
	SignPrivateTesseraTx(context.Context, *big.Int, ethcommon.Address, *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error)

	// SignMsg sign a message
	SignMsg(context.Context, ethcommon.Address, string) ([]byte, *ethcommon.Hash, error)

	// SignRawHash sign a bytes
	SignRawHash(ethcommon.Address, []byte) ([]byte, error)

	// GenerateAccount creates an account
	GenerateAccount(context.Context) (ethcommon.Address, error)

	// ImportPrivateKey creates an account
	ImportPrivateKey(context.Context, string) error
}
