package utils

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
)

// KafkaChainTopic computes Kafka topic identified by chain
func KafkaChainTopic(topic string, chainID *big.Int) string {
	return fmt.Sprintf("%v-%v", topic, chainID.Text(10))
}

var chainAccountKeyPatternRegexp = `(?P<account>0[xX][0-9a-fA-F]{40})@(?P<chain>[0-9]+)`
var chainAccountKeyPattern = regexp.MustCompile(chainAccountKeyPatternRegexp)

// ToChainAccountKey computes a key from a chain identifier and an account
func ToChainAccountKey(chainID *big.Int, acc common.Address) string {
	return fmt.Sprintf("%v@%v", acc.Hex(), chainID.Text(10))
}

// FromChainAccountKey computes a chain identifier and account from a key
func FromChainAccountKey(key string) (chainID *big.Int, acc common.Address, err error) {
	parts := chainAccountKeyPattern.FindStringSubmatch(key)
	if len(parts) != 3 {
		return nil, common.HexToAddress(""), errors.InvalidFormatError("invalid key %q (expected format %q)", key, chainAccountKeyPatternRegexp).SetComponent(component)
	}

	chain, ok := big.NewInt(0).SetString(parts[2], 10)
	if !ok {
		return nil, common.HexToAddress(""), errors.InvalidFormatError("invalid key %q (expected format %q)", key, chainAccountKeyPatternRegexp).SetComponent(component)
	}

	return chain, common.HexToAddress(parts[1]), nil
}
