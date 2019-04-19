package utils

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestKafkaChainTopic(t *testing.T) {
	assert.Equal(t, "test-topic-42", KafkaChainTopic("test-topic", big.NewInt(42)), "Topic should match")
}

func TestChainAccountKey(t *testing.T) {
	assert.Equal(
		t,
		"0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C@3",
		ToChainAccountKey(big.NewInt(3), common.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C")),
		"Key should match",
	)

	chainID, acc, err := FromChainAccountKey("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C@3")

	assert.Nil(t, err, "Should not error")
	assert.Equal(t, int64(3), chainID.Int64(), "ChainID should be correct")
	assert.Equal(t, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", acc.Hex(), "Account should be correct")

	chainID, acc, err = FromChainAccountKey("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C@a3")
	assert.NotNil(t, err, "Should error")
}
