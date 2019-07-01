package utils

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	err "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
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

	chainID, acc, e := FromChainAccountKey("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C@3")
	assert.Nil(t, e, "FromChainAccountKey should not error")
	assert.Equal(t, int64(3), chainID.Int64(), "ChainID should be correct")
	assert.Equal(t, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", acc.Hex(), "Account should be correct")

	_, _, e = FromChainAccountKey("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C@a3")
	assert.NotNil(t, e, "Should error")
	ie, ok := e.(*err.Error)
	assert.True(t, ok, "Error should cast to internal error")
	assert.Equal(t, "utils", ie.GetComponent(), "Error component should be valid")
}
