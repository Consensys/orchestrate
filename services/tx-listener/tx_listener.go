package txlistener

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/session"

	provider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/providers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/session/ethereum"
	hook "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/session/ethereum/hooks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/session/ethereum/offset"
)

type TxListener struct {
	manager session.SManager
}

func NewTxListener(
	prvdr provider.Provider,
	hk hook.Hook,
	offsets offset.Manager,
	ec ethereum.EthClient,
	txSchedulerClient client.TransactionSchedulerClient,
) *TxListener {
	manager := session.NewManager(
		ethereum.NewSessionBuilder(hk, offsets, ec, txSchedulerClient),
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
