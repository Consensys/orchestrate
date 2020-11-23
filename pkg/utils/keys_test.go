// +build unit

package utils

import (
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/error"
)

func TestKafkaChainTopic(t *testing.T) {
	assert.Equal(t, "test-topic-42", KafkaChainTopic("test-topic", big.NewInt(42)), "Topic should match")
}

func TestChainAccountKey(t *testing.T) {
	assert.Equal(
		t,
		"0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C@3",
		ToChainAccountKey(big.NewInt(3), ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C")),
		"Key should match",
	)

	chainID, acc, err := FromChainAddressKey("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C@3")
	assert.NoError(t, err, "FromChainAddressKey should not error")
	assert.Equal(t, int64(3), chainID.Int64(), "GetBigChainID should be correct")
	assert.Equal(t, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", acc.Hex(), "Account should be correct")

	_, _, err = FromChainAddressKey("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C@a3")
	assert.Error(t, err, "Should error")
	ierr, ok := err.(*ierror.Error)
	assert.True(t, ok, "Error should cast to internal error")
	assert.Equal(t, "utils", ierr.GetComponent(), "Error component should be valid")
}
