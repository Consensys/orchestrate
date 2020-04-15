package hook

import (
	"context"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

//go:generate mockgen -source=hook.go -destination=mock/mock.go -package=mock

type Hook interface {
	AfterNewBlock(ctx context.Context, chain *dynamic.Chain, block *ethtypes.Block, receipts []*types.Receipt) error
}
