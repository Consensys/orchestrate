package txlistener

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	provider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/providers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum"
	hook "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/hooks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset"
)

type TxListener struct {
	manager *session.Manager
}

func NewTxListener(prvdr provider.Provider, hk hook.Hook, offsets offset.Manager) *TxListener {
	manager := session.NewManager(
		ethereum.NewSessionBuilder(hk, offsets),
		prvdr,
	)
	return &TxListener{
		manager: manager,
	}
}

func (l *TxListener) Start(ctx context.Context) {
	logger := log.FromContext(ctx)
	l.manager.Start(ctx)
	logger.Infof("Shutting down")
}
