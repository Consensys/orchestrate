package chainregistry

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

const component = "tx-listener.offset"

type Manager struct {
	sm       *sync.Map
	registry registry.Client
}

func NewManager(r registry.Client) *Manager {
	return &Manager{
		sm:       &sync.Map{},
		registry: r,
	}
}

func (m *Manager) GetLastBlockNumber(ctx context.Context, node *dynamic.Node) (int64, error) {
	n, err := m.registry.GetNodeByID(ctx, node.ID)
	if err != nil {
		return 0, errors.FromError(err).ExtendComponent(component)
	}
	return n.ListenerBlockPosition, nil
}

func (m *Manager) SetLastBlockNumber(ctx context.Context, node *dynamic.Node, blockNumber int64) error {
	err := m.registry.UpdateBlockPosition(ctx, node.ID, blockNumber)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}
	return nil
}

func (m *Manager) GetLastTxIndex(_ context.Context, node *dynamic.Node, blockNumber int64) (uint64, error) {
	txIndex, ok := m.sm.Load(fmt.Sprintf("txIndex-%v-%v", node.ID, blockNumber))
	if !ok {
		return 0, nil
	}
	return txIndex.(uint64), nil
}

func (m *Manager) SetLastTxIndex(_ context.Context, node *dynamic.Node, blockNumber int64, txIndex uint64) error {
	m.sm.Store(fmt.Sprintf("txIndex-%v-%v", node.ID, blockNumber), txIndex)
	return nil
}
