package main

import (
	"fmt"
	"math/big"
)

func chainTopic(outTopic string, chainID *big.Int) string {
	return fmt.Sprintf("%v-%v", outTopic, chainID.Text(16))
}
