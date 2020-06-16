package testutils

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"

	"github.com/gofrs/uuid"
)

func FakeChain() *types.Chain {
	return &types.Chain{
		UUID:     uuid.Must(uuid.NewV4()).String(),
		Name:     "FakeChainName",
		TenantID: "_",
	}
}
