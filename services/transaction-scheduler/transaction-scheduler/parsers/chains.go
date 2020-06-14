package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

func NewChainFromModels(chainModel *models.Chain) *types.Chain {
	return &types.Chain{
		Name:     chainModel.Name,
		UUID:     chainModel.UUID,
		TenantID: chainModel.TenantID,
	}
}
