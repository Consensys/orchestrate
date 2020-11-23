package hook

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/dynamic"
)

//go:generate mockgen -source=hook.go -destination=mock/mock.go -package=mock

type Hook interface {
	AfterNewBlock(ctx context.Context, chain *dynamic.Chain, block *ethtypes.Block, jobs []*entities.Job) error
}
