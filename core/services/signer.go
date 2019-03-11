package services

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/types"
)

// TxSigner is an interface to sign transaction
type TxSigner interface {
	// Sign signs a transaction
	Sign(chain *types.Chain, a common.Address, tx *ethtypes.Transaction) (raw []byte, hash *common.Hash, err error)
}
