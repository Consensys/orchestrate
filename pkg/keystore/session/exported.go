package session

import (
	"context"
	"math/big"

	quorumtypes "github.com/consensys/quorum/core/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/account"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/crypto"
)

//go:generate mockgen -source=exported.go -destination=../mock/session.go -package=mock

type SigningSession interface {
	SetAccount(account account.Account) error
	SetChain(*big.Int) error
	ExecuteForTx(*ethtypes.Transaction) ([]byte, *ethcommon.Hash, error)
	ExecuteForMsg([]byte, crypto.DSA) ([]byte, *ethcommon.Hash, error)
	ExecuteForTesseraTx(*quorumtypes.Transaction) ([]byte, *ethcommon.Hash, error)
	ExecuteForEEATx(*ethtypes.Transaction, *types.PrivateArgs) ([]byte, *ethcommon.Hash, error)
}

type AccountManager interface {
	SigningSession(context.Context, ethcommon.Address) (SigningSession, error)
	ImportPrivateKey(context.Context, string) error
	GenerateAccount(context.Context) (ethcommon.Address, error)
}
