package listener

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func TestTxListener(t *testing.T) {
	blocks := []*types.Block{}
	for _, blockEnc := range blocksEnc {
		var block types.Block
		err := rlp.DecodeBytes(blockEnc, &block)
		assert.Nil(t, err)
		blocks = append(blocks, &block)
	}
	mec := NewMockEthClient(blocks)

	// Initialize Configuration
	config := NewConfig()
	config.BlockCursor.Backoff = 100 * time.Millisecond
	config.BlockCursor.Limit = uint64(40)
	config.BlockCursor.Tracker.Depth = uint64(0)
	config.TxListener.Return.Blocks = false
	config.TxListener.Return.Errors = false

	l := NewTxListener(mec, config)
	_, err := l.Listen(big.NewInt(1), 0, 0)
	if err != nil {
		t.Errorf("TxListener #1: expected no error but got %v", err)
	}

	// Try to listen to the same chain again
	_, err = l.Listen(big.NewInt(1), 0, 0)
	if err == nil {
		t.Errorf("TxListener #2: expected an error")
	}

	// Try to listen to another chain
	_, err = l.Listen(big.NewInt(2), 1, 5)
	if err != nil {
		t.Errorf("TxListener #3: expected no error but got %v", err)
	}

	wait := &sync.WaitGroup{}
	wait.Add(3)

	blcks := []*TxListenerBlock{}
	go func() {
		for block := range l.Blocks() {
			// Drain blocks
			blcks = append(blcks, block)
		}
		wait.Done()
	}()

	receipts := []*TxListenerReceipt{}
	go func() {
		for receipt := range l.Receipts() {
			receipts = append(receipts, receipt)
		}
		wait.Done()
	}()

	errors := []error{}
	go func() {
		for err := range l.Errors() {
			// Drain blocks
			errors = append(errors, err)
		}
		wait.Done()
	}()

	// Simulate 2 mined blocks
	time.Sleep(50 * time.Millisecond)
	mec.mine()

	// Try to listen to another chain
	_, err = l.Listen(big.NewInt(3), -1, 0)
	if err != nil {
		t.Errorf("TxListener #4: expected no error but got %v", err)
	}

	time.Sleep(180 * time.Millisecond)
	mec.mine()

	time.Sleep(10 * time.Millisecond)
	mec.mine()

	// Test methods while running (to maybe detect some race conditions)
	chains := l.Chains()
	expected := 3
	if len(chains) != expected {
		t.Errorf("TxListener: expected %v chain but got %v", expected, len(chains))
	}

	progress := l.Progress(context.Background())
	highest := 3
	if progress["1"].HighestBlock != int64(highest) {
		t.Errorf("TxListener: expected highest block to be %v but got %v", highest, progress["1"].HighestBlock)
	}

	// Close listener
	l.Close()

	// Try to listen on a close listener
	_, err = l.Listen(big.NewInt(1), 0, 0)
	if err == nil {
		t.Errorf("TxListener #2: expected an error")
	}

	// Wait for drainers to complete
	wait.Wait()

	// Test methods after stoping
	chains = l.Chains()
	expected = 0
	if len(chains) != expected {
		t.Errorf("TxListener: expected %v chain but got %v", expected, len(chains))
	}

	if len(errors) != 0 {
		t.Errorf("TxListener: expected %v errors but got %v", 0, len(errors))
	}

	if len(blcks) != 0 {
		t.Errorf("TxListener: expected %v blocks but got %v", 0, len(blcks))
	}

	expected = 31
	if len(receipts) != expected {
		t.Errorf("TxListener: expected %v receipts but got %v", expected, len(receipts))
	}
}

func TestTxListenerWithReturns(t *testing.T) {
	blocks := []*types.Block{}
	for _, blockEnc := range blocksEnc {
		var block types.Block
		err := rlp.DecodeBytes(blockEnc, &block)
		assert.Nil(t, err)
		blocks = append(blocks, &block)
	}
	mec := NewMockEthClient(blocks)

	// Initialize cursor
	config := NewConfig()
	config.BlockCursor.Backoff = 100 * time.Millisecond
	config.BlockCursor.Limit = uint64(40)
	config.BlockCursor.Tracker.Depth = uint64(0)
	config.TxListener.Return.Blocks = true
	config.TxListener.Return.Errors = true

	l := NewTxListener(mec, config)
	_, err := l.Listen(big.NewInt(1), -1, 0)
	if err != nil {
		t.Errorf("TxListener #1: expected no error but got %v", err)
	}

	// Try to listen to another chain
	_, err = l.Listen(big.NewInt(2), 1, 5)
	if err != nil {
		t.Errorf("TxListener #3: expected no error but got %v", err)
	}

	wait := &sync.WaitGroup{}
	wait.Add(3)

	blcks := []*TxListenerBlock{}
	go func() {
		for block := range l.Blocks() {
			// Drain blocks
			blcks = append(blcks, block)
		}
		wait.Done()
	}()

	receipts := []*TxListenerReceipt{}
	go func() {
		for receipt := range l.Receipts() {
			receipts = append(receipts, receipt)
		}
		wait.Done()
	}()

	errors := []error{}
	go func() {
		for err := range l.Errors() {
			// Drain blocks
			errors = append(errors, err)
		}
		wait.Done()
	}()

	// Simulate mined blocks
	time.Sleep(50 * time.Millisecond)
	mec.mine()

	// Try to listen to another chain
	_, err = l.Listen(big.NewInt(3), -1, 0)
	if err != nil {
		t.Errorf("TxListener #4: expected no error but got %v", err)
	}

	// Simulate mined blocks
	time.Sleep(180 * time.Millisecond)
	mec.mine()

	// Simulate mined blocks
	time.Sleep(10 * time.Millisecond)
	mec.mine()

	// Test methods while running
	chains := l.Chains()
	expected := 3
	if len(chains) != expected {
		t.Errorf("TxListener: expected %v chain but got %v", expected, len(chains))
	}

	progress := l.Progress(context.Background())
	highest := 3
	if progress["1"].HighestBlock != int64(highest) {
		t.Errorf("TxListener: expected highest block to be %v but got %v", highest, progress["1"].HighestBlock)
	}

	// Close
	l.Close()

	// Try to listen on a close listener
	_, err = l.Listen(big.NewInt(1), 0, 0)
	if err == nil {
		t.Errorf("TxListener #2: expected an error")
	}

	// Wait for drainers to complete
	wait.Wait()

	if len(errors) != 0 {
		t.Errorf("TxListener: expected %v errors but got %v", 0, len(errors))
	}

	expected = 4
	if len(blcks) != expected {
		t.Errorf("TxListener: expected %v blocks but got %v", expected, len(blcks))
	}

	expected = 31
	if len(receipts) != expected {
		t.Errorf("TxListener: expected %v receipts but got %v", expected, len(receipts))
	}
}
