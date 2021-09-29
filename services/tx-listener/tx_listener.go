package txlistener

import (
	"context"

	orchestrateclient "github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/services/tx-listener/metrics"
	"github.com/consensys/orchestrate/services/tx-listener/session"

	provider "github.com/consensys/orchestrate/services/tx-listener/providers"
	"github.com/consensys/orchestrate/services/tx-listener/session/ethereum"
	hook "github.com/consensys/orchestrate/services/tx-listener/session/ethereum/hooks"
	"github.com/consensys/orchestrate/services/tx-listener/session/ethereum/offset"
)

type TxListener struct {
	manager session.SManager
}

func NewTxListener(
	prvdr provider.Provider,
	hk hook.Hook,
	offsets offset.Manager,
	ec ethereum.EthClient,
	client orchestrateclient.OrchestrateClient,
	m metrics.ListenerMetrics,
) *TxListener {
	manager := session.NewManager(
		ethereum.NewSessionBuilder(hk, offsets, ec, client, m),
		prvdr,
	)

	return &TxListener{
		manager: manager,
	}
}

func (l *TxListener) Run(ctx context.Context) error {
	return l.manager.Run(ctx)
}

func (l *TxListener) Close() error {
	return nil
}
