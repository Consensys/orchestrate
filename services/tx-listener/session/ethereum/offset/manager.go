package offset

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

//go:generate mockgen -source=manager.go -destination=mock/mock.go -package=mock

type Manager interface {
	GetLastBlockNumber(ctx context.Context, chain *dynamic.Chain) (uint64, error)
	SetLastBlockNumber(ctx context.Context, chain *dynamic.Chain, blockNumber uint64) error
	GetLastTxIndex(ctx context.Context, chain *dynamic.Chain, blockNumber uint64) (uint64, error)
	SetLastTxIndex(ctx context.Context, chain *dynamic.Chain, blockNumber uint64, txIndex uint64) error
}
