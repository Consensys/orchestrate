// +build unit

package ethereum

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	// Test from GethLog
	ethlog := &ethtypes.Log{
		Address: ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"),
		Topics: []ethcommon.Hash{
			ethcommon.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
			ethcommon.HexToHash("0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615"),
		},
		Data:        []byte{0x1a, 0x05, 0x56, 0x90, 0xd9, 0xdb, 0x80, 0x00},
		BlockNumber: 2019236,
		TxHash:      ethcommon.HexToHash("0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e"),
		TxIndex:     3,
		BlockHash:   ethcommon.HexToHash("0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056"),
		Removed:     true,
	}

	log := FromGethLog(ethlog)
	assert.Equal(t, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", log.GetAddress(), "Address should match")
	var topics []string
	topics = append(topics, log.GetTopics()...)
	assert.Equal(t, []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615"}, topics, "Topics should match")
	assert.Equal(t, hexutil.Encode([]byte{0x1a, 0x05, 0x56, 0x90, 0xd9, 0xdb, 0x80, 0x00}), log.Data, "Data should match")
	assert.Equal(t, uint64(2019236), log.BlockNumber, "Blocknumber should match")
	assert.Equal(t, "0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e", log.GetTxHash(), "TxHash should match")
	assert.Equal(t, uint64(3), log.TxIndex, "TxIndex should match")
	assert.Equal(t, "0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056", log.GetBlockHash(), "BlockHash should match")
	assert.True(t, log.Removed, "Removed should match")

	var topicHashes []ethcommon.Hash
	for _, topic := range log.GetTopics() {
		topicHashes = append(topicHashes, ethcommon.HexToHash(topic))
	}
	assert.Equal(t, ethlog.Topics, topicHashes, "Topics should match")
}

func TestReceipt(t *testing.T) {
	gethReceipt := &ethtypes.Receipt{
		ContractAddress:   ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"),
		PostState:         []byte{0x3b, 0x19, 0x8b, 0xfd, 0x5d, 0x29, 0x07, 0x28, 0x5a, 0xf0, 0x09, 0xe9, 0xae, 0x84, 0xa0, 0xec, 0xd6, 0x36, 0x77, 0x11, 0x0d, 0x89, 0xd7, 0xe0, 0x30, 0x25, 0x1a, 0xcb, 0x87, 0xf6, 0x48, 0x7e},
		Status:            1,
		TxHash:            ethcommon.HexToHash("0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e"),
		Bloom:             ethtypes.BytesToBloom([]byte{0x1a, 0x05, 0x56, 0x90, 0xd9, 0xdb, 0x80, 0x00}),
		GasUsed:           uint64(156),
		CumulativeGasUsed: uint64(14567),
		Logs: []*ethtypes.Log{
			{
				Address: ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"),
				Topics: []ethcommon.Hash{
					ethcommon.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
					ethcommon.HexToHash("0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615"),
				},
				Data:        []byte{0x1a, 0x05, 0x56, 0x90, 0xd9, 0xdb, 0x80, 0x00},
				BlockNumber: 2019236,
				TxHash:      ethcommon.HexToHash("0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e"),
				TxIndex:     3,
				BlockHash:   ethcommon.HexToHash("0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056"),
				Removed:     true,
			},
		},
	}

	r := FromGethReceipt(gethReceipt)
	assert.Equal(t, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", r.GetContractAddress(), "ContractAddress should match")
	assert.Equal(t, hexutil.Encode(gethReceipt.PostState), r.PostState, "PostState should match")
	assert.Equal(t, uint64(1), r.Status, "Status should match")
	assert.Equal(t, "0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e", r.GetTxHash(), "TxHash should match")
	assert.Equal(t, hexutil.Encode(gethReceipt.Bloom.Bytes()), r.Bloom, "Bloom should match")
	assert.Equal(t, uint64(156), r.GasUsed, "GasUsed should match")
	assert.Equal(t, uint64(14567), r.CumulativeGasUsed, "CumulativeGasUsed should match")

	assert.Len(t, r.Logs, 1, "Logs should be 1 long")
	assert.Equal(t, uint64(2019236), r.Logs[0].BlockNumber, "Logs should match")

	// Test Setters
	r = r.
		SetBlockHash(ethcommon.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")).
		SetBlockNumber(123456).
		SetTxIndex(123)

	assert.Equal(t, "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", r.GetBlockHash(), "BlockHash should match")
	assert.Equal(t, uint64(123456), r.BlockNumber, "BlockNumber should match")
	assert.Equal(t, uint64(123), r.TxIndex, "TxIndex should match")
}
