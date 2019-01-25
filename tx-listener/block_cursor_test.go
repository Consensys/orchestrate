package listener

import (
	"math/big"
	"reflect"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestBlockCursor(t *testing.T) {

	// Create a cursor
	cursor := &blockCursor{
		c: &mockEthClient{
			t:        t,
			mux:      &sync.Mutex{},
			index:    0,
			expected: []string{"mined"},
		},
		next: big.NewInt(0),
	}

	block, err := cursor.Next()
	if err != nil {
		t.Errorf("Next: expected to get block but got %v", err)
	}

	if !reflect.DeepEqual(block.Hash(), common.HexToHash("80bd976d96ef1da0251150e741fd596d3e580be70b02a4757554a452c17edfe1")) {
		t.Errorf("Next: Hash mismatch got %v, want %v", block.Hash().Hex(), "0x80bd976d96ef1da0251150e741fd596d3e580be70b02a4757554a452c17edfe1")
	}

	if cursor.next.Uint64() != 1 {
		t.Errorf("Next: Cursor should have been incremented")
	}

	block, err = cursor.Next()
	if err == nil || block != nil {
		t.Errorf("Next: expected an error")
	}

	if cursor.next.Uint64() != 1 {
		t.Errorf("Next: Cursor should not have been incremented")
	}

	cursor.Set(big.NewInt(10))
	if cursor.next.Uint64() != 10 {
		t.Errorf("Next: Cursor should have been set to %v but got %v", 10, cursor.next.Uint64())
	}
}
