package offset

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Manager interface {
	GetLastBlockNumber(ctx context.Context, chain *dynamic.Chain) (int64, error)
	SetLastBlockNumber(ctx context.Context, chain *dynamic.Chain, blockNumber int64) error
	GetLastTxIndex(ctx context.Context, chain *dynamic.Chain, blockNumber int64) (uint64, error)
	SetLastTxIndex(ctx context.Context, chain *dynamic.Chain, blockNumber int64, txIndex uint64) error
}
