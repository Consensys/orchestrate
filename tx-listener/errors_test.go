package listener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReceiptMissingError(t *testing.T) {
	err := &ReceiptMissingError{"0xhash"}
	assert.Equal(t, "Receipt for Transaction \"0xhash\" missing", err.Error(), "Error message should match")
}

func TestBlockMissingError(t *testing.T) {
	err := &BlockMissingError{int64(43)}
	assert.Equal(t, "Block 43 missing", err.Error(), "Error message should match")
}
