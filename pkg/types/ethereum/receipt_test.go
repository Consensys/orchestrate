// +build unit

package ethereum

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

func Test_unpackRevertReason(t *testing.T) {
	message, err := unpackRevertReason(ethcommon.FromHex("0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002645524332303a207472616e7366657220616d6f756e7420657863656564732062616c616e63650000000000000000000000000000000000000000000000000000"))
	assert.NoError(t, err)
	assert.Equal(t, "ERC20: transfer amount exceeds balance", message)
}

func TestReceipt_UnmarshalJSON(t *testing.T) {
	type fields Receipt
	type args struct {
		input []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"",
			fields{
				TxHash:            "0x8c16a01f417822f19057e0d64a87497956267b0017b65a5364fa38aa7af83555",
				BlockHash:         "0x06efb4ac4e192b8c8a4a4d8ad7c009168470c5bc1c36cf116689f71ea36f08e3",
				BlockNumber:       36,
				TxIndex:           0,
				ContractAddress:   "",
				PostState:         "",
				Status:            0,
				Bloom:             "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				GasUsed:           24103,
				CumulativeGasUsed: 24103,
				RevertReason:      "ERC20: transfer amount exceeds balance",
			},
			args{
				[]byte(`{
    "blockHash" : "0x06efb4ac4e192b8c8a4a4d8ad7c009168470c5bc1c36cf116689f71ea36f08e3",
    "blockNumber" : "0x24",
    "contractAddress" : null,
    "cumulativeGasUsed" : "0x5e27",
    "from" : "0x6009608a02a7a15fd6689d6dad560c44e9ab61ff",
    "gasUsed" : "0x5e27",
    "logs" : [ ],
    "logsBloom" : "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "status" : "0x0",
    "to" : "0x71b7d704598945e72e7581bac3b070d300dc6eb3",
    "transactionHash" : "0x8c16a01f417822f19057e0d64a87497956267b0017b65a5364fa38aa7af83555",
    "transactionIndex" : "0x0",
    "revertReason" : "0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002645524332303a207472616e7366657220616d6f756e7420657863656564732062616c616e63650000000000000000000000000000000000000000000000000000"
}`),
			},
			false,
		},
		{
			"",
			fields{
			},
			args{
				[]byte(`{
    "transactionHash" : "0xinvalidHash",
}`),
			},
			true,
		},
		{
			"",
			fields{
			},
			args{
				[]byte(`{`),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := &Receipt{
				TxHash:               tt.fields.TxHash,
				BlockHash:            tt.fields.BlockHash,
				BlockNumber:          tt.fields.BlockNumber,
				TxIndex:              tt.fields.TxIndex,
				ContractAddress:      tt.fields.ContractAddress,
				PostState:            tt.fields.PostState,
				Status:               tt.fields.Status,
				Bloom:                tt.fields.Bloom,
				Logs:                 tt.fields.Logs,
				GasUsed:              tt.fields.GasUsed,
				CumulativeGasUsed:    tt.fields.CumulativeGasUsed,
				RevertReason:         tt.fields.RevertReason,
				XXX_NoUnkeyedLiteral: tt.fields.XXX_NoUnkeyedLiteral,
				XXX_unrecognized:     tt.fields.XXX_unrecognized,
				XXX_sizecache:        tt.fields.XXX_sizecache,
			}
			r := &Receipt{}
			if err := r.UnmarshalJSON(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, expected, r)
		})
	}
}
