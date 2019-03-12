package services

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

// TxSigner is an interface to sign transaction
type TxSigner interface {
	// Sign signs a transaction
	Sign(chain *common.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, hash *ethcommon.Hash, err error)
}
