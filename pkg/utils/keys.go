package utils

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/consensys/orchestrate/pkg/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// KafkaChainTopic computes Kafka topic identified by chain
func KafkaChainTopic(topic string, chainID *big.Int) string {
	return fmt.Sprintf("%v-%v", topic, chainID.Text(10))
}

var chainAddressKeyPatternRegexp = `(?P<address>0[xX][0-9a-fA-F]{40})@(?P<chain>[0-9]+)`
var chainAddressKeyPattern = regexp.MustCompile(chainAddressKeyPatternRegexp)

// ToChainAccountKey computes a key from a chain identifier and an account
func ToChainAccountKey(chainID *big.Int, acc ethcommon.Address) string {
	return fmt.Sprintf("%v@%v", acc.Hex(), chainID.Text(10))
}

// FromChainAddressKey computes a chain identifier and account from a key
func FromChainAddressKey(key string) (chainID *big.Int, acc ethcommon.Address, err error) {
	parts := chainAddressKeyPattern.FindStringSubmatch(key)
	if len(parts) != 3 {
		return nil, ethcommon.HexToAddress(""), errors.InvalidFormatError("invalid key %q (expected format %q)", key, chainAddressKeyPatternRegexp).SetComponent(component)
	}

	chain, ok := big.NewInt(0).SetString(parts[2], 10)
	if !ok {
		return nil, ethcommon.HexToAddress(""), errors.InvalidFormatError("invalid key %q (expected format %q)", key, chainAddressKeyPatternRegexp).SetComponent(component)
	}

	return chain, ethcommon.HexToAddress(parts[1]), nil
}
