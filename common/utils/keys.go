package utils

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
)

// KafkaChainTopic computes Kafka topic identified by chain
func KafkaChainTopic(topic string, chainID *big.Int) string {
	return fmt.Sprintf("%v-%v", topic, chainID.Text(16))
}

var chainAccountKeyPatternRegexp = `(?P<account>\w+)@(?P<chain>\w+)`
var chainAccountKeyPattern = regexp.MustCompile(chainAccountKeyPatternRegexp)

// ToChainAccountKey computes a key from a chain identifier and an account
func ToChainAccountKey(chainID *big.Int, acc common.Address) string {
	return fmt.Sprintf("%v@%v", acc.Hex(), chainID.Text(16))
}

// FromChainAccountKey computes a chain identifier and account from a key
func FromChainAccountKey(key string) (chainID *big.Int, acc common.Address, err error) {
	parts := chainAccountKeyPattern.FindStringSubmatch(key)
	if len(parts) != 3 {
		return nil, common.HexToAddress(""), fmt.Errorf("ChainAccountKey %q invalid (expect %q)", key, chainAccountKeyPatternRegexp)
	}

	chain, ok := big.NewInt(0).SetString(parts[2], 16)
	if !ok {
		return nil, common.HexToAddress(""), fmt.Errorf("ChainID %q is an invalid hexadecimal", parts[2])
	}

	return chain, common.HexToAddress(parts[1]), nil
}
