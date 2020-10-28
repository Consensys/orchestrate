package utils

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

func ParseHexToMixedCaseEthAddress(address string) (string, error) {
	if !ethcommon.IsHexAddress(address) {
		return "", fmt.Errorf("expected hex string")
	}

	addr := ethcommon.HexToAddress(address)
	return addr.String(), nil
}
