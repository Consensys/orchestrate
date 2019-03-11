package ethereum

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	// Test from GethLog
	gethLog := &types.Log{
		Address: common.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"),
		Topics: []common.Hash{
			common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
			common.HexToHash("0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615"),
		},
		Data:        []byte{0x1a, 0x05, 0x56, 0x90, 0xd9, 0xdb, 0x80, 0x00},
		BlockNumber: 2019236,
		TxHash:      common.HexToHash("0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e"),
		TxIndex:     3,
		BlockHash:   common.HexToHash("0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056"),
		Removed:     true,
	}

	log := FromGethLog(gethLog)
	assert.Equal(t, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", log.Address, "Address should match")
	assert.Equal(t, []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615"}, log.Topics, "Topics should match")
	assert.Equal(t, []byte{0x1a, 0x05, 0x56, 0x90, 0xd9, 0xdb, 0x80, 0x00}, log.Data, "Data should match")
	assert.Equal(t, uint64(2019236), log.BlockNumber, "Blocknumber should match")
	assert.Equal(t, "0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e", log.TxHash, "TxHash should match")
	assert.Equal(t, uint64(3), log.TxIndex, "TxIndex should match")
	assert.Equal(t, "0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056", log.BlockHash, "BlockHash should match")
	assert.Equal(t, true, log.Removed, "Removed should match")

	assert.Equal(t, gethLog.Topics, log.TopicsHash(), "TopicsHash should match")
}

func TestReceipt(t *testing.T) {
	gethReceipt := &types.Receipt{
		ContractAddress:   common.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"),
		PostState:         []byte{0x3b, 0x19, 0x8b, 0xfd, 0x5d, 0x29, 0x07, 0x28, 0x5a, 0xf0, 0x09, 0xe9, 0xae, 0x84, 0xa0, 0xec, 0xd6, 0x36, 0x77, 0x11, 0x0d, 0x89, 0xd7, 0xe0, 0x30, 0x25, 0x1a, 0xcb, 0x87, 0xf6, 0x48, 0x7e},
		Status:            1,
		TxHash:            common.HexToHash("0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e"),
		Bloom:             types.BytesToBloom([]byte{0x1a, 0x05, 0x56, 0x90, 0xd9, 0xdb, 0x80, 0x00}),
		GasUsed:           uint64(156),
		CumulativeGasUsed: uint64(14567),
		Logs: []*types.Log{
			&types.Log{
				Address: common.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"),
				Topics: []common.Hash{
					common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
					common.HexToHash("0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615"),
				},
				Data:        []byte{0x1a, 0x05, 0x56, 0x90, 0xd9, 0xdb, 0x80, 0x00},
				BlockNumber: 2019236,
				TxHash:      common.HexToHash("0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e"),
				TxIndex:     3,
				BlockHash:   common.HexToHash("0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056"),
				Removed:     true,
			},
		},
	}

	r := FromGethReceipt(gethReceipt)
	assert.Equal(t, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", r.ContractAddress, "ContractAddress should match")
	assert.Equal(t, gethReceipt.PostState, r.PostState, "PostState should match")
	assert.Equal(t, uint64(1), r.Status, "Status should match")
	assert.Equal(t, "0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e", r.TxHash, "TxHash should match")
	assert.Equal(t, gethReceipt.Bloom.Bytes(), r.Bloom, "Bloom should match")
	assert.Equal(t, uint64(156), r.GasUsed, "GasUsed should match")
	assert.Equal(t, uint64(14567), r.CumulativeGasUsed, "CumulativeGasUsed should match")

	assert.Len(t, r.Logs, 1, "Logs should be 1 long")
	assert.Equal(t, uint64(2019236), r.Logs[0].BlockNumber, "Logs should match")

	// Test Setters
	r = r.
		SetBlockHash(common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")).
		SetBlockNumber(123456).
		SetTxIndex(123)

	assert.Equal(t, "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", r.BlockHash, "BlockHash should match")
	assert.Equal(t, uint64(123456), r.BlockNumber, "BlockNumber should match")
	assert.Equal(t, uint64(123), r.TxIndex, "TxIndex should match")
}
