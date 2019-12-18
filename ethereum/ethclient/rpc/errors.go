package rpc

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

var (
	EthConnErr = errors.EthConnectionError("could not connect to Ethereum client")
)
