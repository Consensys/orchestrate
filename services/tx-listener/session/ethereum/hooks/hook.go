package hook

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

//go:generate mockgen -source=hook.go -destination=mock/mock.go -package=mock

type Hook interface {
	AfterNewBlockEnvelope(ctx context.Context, chain *dynamic.Chain, block *ethtypes.Block, envelopes []*tx.Envelope) error
	AfterNewBlock(ctx context.Context, chain *dynamic.Chain, block *ethtypes.Block, jobs []*entities.Job) error
}
