package listener

import (
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestBlockListener(t *testing.T) {
	// Create mock client
	ec := &mockEthClient{
		t:        t,
		mux:      &sync.Mutex{},
		index:    0,
		expected: []string{"mined", "mined", "pending", "pending", "mined", "pending", "mined", "mined", "pending", "mined"},
	}

	// Create cursor & block consumer
	cursor := &blockCursor{
		c:    ec,
		next: big.NewInt(0),
	}
	config := NewConfig()
	config.BlockListener.Backoff = 10 * time.Millisecond
	config.BlockListener.Return.Errors = true
	bl := newBlockListener(cursor, config)

	// Run feeder
	go bl.feeder()

	go func() {
		// Close on error
		<-bl.Errors()
		bl.Close()
	}()

	blocks := []*types.Block{}
	// Drain blocks
	for block := range bl.Blocks() {
		blocks = append(blocks, block)
	}

	if len(blocks) != 6 {
		t.Errorf("BlockConsumer: expected 6 blocks but got %v", len(blocks))
	}
}

func TestBlockListenerInterupted(t *testing.T) {
	// Create mock client
	ec := &mockEthClient{
		t:        t,
		mux:      &sync.Mutex{},
		index:    0,
		expected: []string{"mined", "mined", "pending", "pending", "mined", "pending", "mined", "mined", "pending", "mined"},
	}

	// Create cursor & block consumer
	cursor := &blockCursor{
		c:    ec,
		next: big.NewInt(0),
	}
	config := NewConfig()
	config.BlockListener.Backoff = 10 * time.Millisecond
	bl := newBlockListener(cursor, config)

	// Run feeder
	go bl.feeder()

	// Simulate a close before first error
	go func() {
		time.Sleep(25 * time.Millisecond)
		bl.Close()
	}()

	blocks := []*types.Block{}
	// Drain blocks
	for block := range bl.Blocks() {
		blocks = append(blocks, block)
	}

	if len(blocks) != 3 {
		t.Errorf("BlockConsumer: expected %v blocks but got %v", 3, len(blocks))
	}
}
