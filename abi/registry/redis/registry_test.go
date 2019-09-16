package redis

import (
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry/suites"
)

func TestRegistry(t *testing.T) {
	suites.RunRegistry(t, NewRegistry(NewPool(Config(), DialMock)))
}