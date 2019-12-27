package memory

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Manager struct {
	mux   *sync.Mutex
	cache map[string]uint64
}

func NewManager() *Manager {
	return &Manager{
		mux:   &sync.Mutex{},
		cache: make(map[string]uint64),
	}
}

func (m *Manager) GetLastBlockNumber(ctx context.Context, node *dynamic.Node) (uint64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	blockNumber, ok := m.cache[fmt.Sprintf("blockNumber-%v", node.ID)]
	if !ok {
		return 0, nil
	}
	return blockNumber, nil
}

func (m *Manager) SetLastBlockNumber(ctx context.Context, node *dynamic.Node, blockNumber uint64) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.cache[fmt.Sprintf("blockNumber-%v", node.ID)] = blockNumber
	return nil
}

func (m *Manager) GetLastTxIndex(ctx context.Context, node *dynamic.Node, blockNumber uint64) (uint64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	txIndex, ok := m.cache[fmt.Sprintf("txIndex-%v-%v", node.ID, blockNumber)]
	if !ok {
		return 0, nil
	}
	return txIndex, nil
}

func (m *Manager) SetLastTxIndex(ctx context.Context, node *dynamic.Node, blockNumber, txIndex uint64) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.cache[fmt.Sprintf("txIndex-%v-%v", node.ID, blockNumber)] = txIndex
	return nil
}
