package memory

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Manager struct {
	mux   *sync.Mutex
	cache map[string]int64
}

func NewManager() *Manager {
	return &Manager{
		mux:   &sync.Mutex{},
		cache: make(map[string]int64),
	}
}

func (m *Manager) GetLastBlockNumber(_ context.Context, node *dynamic.Node) (int64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	blockNumber, ok := m.cache[fmt.Sprintf("blockNumber-%v", node.ID)]
	if !ok {
		return 0, nil
	}
	return blockNumber, nil
}

func (m *Manager) SetLastBlockNumber(_ context.Context, node *dynamic.Node, blockNumber int64) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.cache[fmt.Sprintf("blockNumber-%v", node.ID)] = blockNumber
	return nil
}

func (m *Manager) GetLastTxIndex(_ context.Context, node *dynamic.Node, blockNumber int64) (uint64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	txIndex, ok := m.cache[fmt.Sprintf("txIndex-%v-%v", node.ID, blockNumber)]
	if !ok {
		return 0, nil
	}
	return uint64(txIndex), nil
}

func (m *Manager) SetLastTxIndex(_ context.Context, node *dynamic.Node, blockNumber int64, txIndex uint64) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.cache[fmt.Sprintf("txIndex-%v-%v", node.ID, blockNumber)] = int64(txIndex)
	return nil
}
