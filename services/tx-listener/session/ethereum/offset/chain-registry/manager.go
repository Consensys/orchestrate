package chainregistry

import (
	"context"
	"fmt"
	"sync"

	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/dynamic"
)

const component = "tx-listener.offset"

type Manager struct {
	sm     *sync.Map
	client orchestrateclient.ChainClient
}

func NewManager(client orchestrateclient.ChainClient) *Manager {
	return &Manager{
		sm:     &sync.Map{},
		client: client,
	}
}

func (m *Manager) GetLastBlockNumber(ctx context.Context, chain *dynamic.Chain) (uint64, error) {
	chainRetrieved, err := m.client.GetChain(ctx, chain.UUID)
	if err != nil {
		return 0, errors.FromError(err).ExtendComponent(component)
	}

	return chainRetrieved.ListenerCurrentBlock, nil
}

func (m *Manager) SetLastBlockNumber(ctx context.Context, chain *dynamic.Chain, blockNumber uint64) error {
	_, err := m.client.UpdateChain(ctx, chain.UUID, &api.UpdateChainRequest{Listener: &api.UpdateListenerRequest{CurrentBlock: blockNumber}})
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}
	return nil
}

func (m *Manager) GetLastTxIndex(_ context.Context, chain *dynamic.Chain, blockNumber uint64) (uint64, error) {
	txIndex, ok := m.sm.Load(fmt.Sprintf("txIndex-%v-%v", chain.UUID, blockNumber))
	if !ok {
		return 0, nil
	}
	return txIndex.(uint64), nil
}

func (m *Manager) SetLastTxIndex(_ context.Context, chain *dynamic.Chain, blockNumber, txIndex uint64) error {
	m.sm.Store(fmt.Sprintf("txIndex-%v-%v", chain.UUID, blockNumber), txIndex)
	return nil
}
