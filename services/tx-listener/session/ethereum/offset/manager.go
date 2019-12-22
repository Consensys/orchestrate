package offset

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Manager interface {
	GetLastBlockNumber(ctx context.Context, node *dynamic.Node) (uint64, error)
	SetLastBlockNumber(ctx context.Context, node *dynamic.Node, blockNumber uint64) error
	GetLastTxIndex(ctx context.Context, node *dynamic.Node, blockNumder uint64) (uint64, error)
	SetLastTxIndex(ctx context.Context, node *dynamic.Node, blockNumber, txIndex uint64) error
}
