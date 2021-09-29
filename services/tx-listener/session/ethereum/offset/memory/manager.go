package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/consensys/orchestrate/services/tx-listener/dynamic"
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

func (m *Manager) GetLastBlockNumber(_ context.Context, chain *dynamic.Chain) (uint64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	blockNumber, ok := m.cache[fmt.Sprintf("blockNumber-%v", chain.UUID)]
	if !ok {
		return 0, nil
	}
	return blockNumber, nil
}

func (m *Manager) SetLastBlockNumber(_ context.Context, chain *dynamic.Chain, blockNumber uint64) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.cache[fmt.Sprintf("blockNumber-%v", chain.UUID)] = blockNumber
	return nil
}

func (m *Manager) GetLastTxIndex(_ context.Context, chain *dynamic.Chain, blockNumber uint64) (uint64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	txIndex, ok := m.cache[fmt.Sprintf("txIndex-%v-%v", chain.UUID, blockNumber)]
	if !ok {
		return 0, nil
	}
	return txIndex, nil
}

func (m *Manager) SetLastTxIndex(_ context.Context, chain *dynamic.Chain, blockNumber, txIndex uint64) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.cache[fmt.Sprintf("txIndex-%v-%v", chain.UUID, blockNumber)] = txIndex
	return nil
}
