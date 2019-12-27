package hook

import (
	"context"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Hook interface {
	AfterNewBlock(ctx context.Context, node *dynamic.Node, block *ethtypes.Block, receipts []*ethtypes.Receipt) error
}
