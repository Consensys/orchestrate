package parsers

import (
	"math/big"

	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/services/api/store/models"
)

func NewChainFromModel(chainModel *models.Chain) *entities.Chain {
	chain := &entities.Chain{
		UUID:                      chainModel.UUID,
		Name:                      chainModel.Name,
		TenantID:                  chainModel.TenantID,
		OwnerID:                   chainModel.OwnerID,
		URLs:                      chainModel.URLs,
		ChainID:                   (*big.Int)(utils.BigIntStringToHex(chainModel.ChainID)),
		ListenerDepth:             chainModel.ListenerDepth,
		ListenerCurrentBlock:      chainModel.ListenerCurrentBlock,
		ListenerStartingBlock:     chainModel.ListenerStartingBlock,
		ListenerBackOffDuration:   chainModel.ListenerBackOffDuration,
		ListenerExternalTxEnabled: chainModel.ListenerExternalTxEnabled,
		Labels:                    chainModel.Labels,
		Headers:                   chainModel.Headers,
		CreatedAt:                 chainModel.CreatedAt,
		UpdatedAt:                 chainModel.UpdatedAt,
	}

	if len(chainModel.PrivateTxManagers) > 0 {
		chain.PrivateTxManager = NewPrivateTxManagerFromModel(chainModel.PrivateTxManagers[0])
	}

	return chain
}

func NewPrivateTxManagerFromModel(privateTxManager *models.PrivateTxManager) *entities.PrivateTxManager {
	return &entities.PrivateTxManager{
		UUID:      privateTxManager.UUID,
		ChainUUID: privateTxManager.ChainUUID,
		URL:       privateTxManager.URL,
		Type:      privateTxManager.Type,
		CreatedAt: privateTxManager.CreatedAt,
	}
}

func NewChainModelFromEntity(chain *entities.Chain) *models.Chain {
	chainModel := &models.Chain{
		UUID:                      chain.UUID,
		Name:                      chain.Name,
		TenantID:                  chain.TenantID,
		OwnerID:                   chain.OwnerID,
		URLs:                      chain.URLs,
		ListenerDepth:             chain.ListenerDepth,
		ListenerCurrentBlock:      chain.ListenerCurrentBlock,
		ListenerStartingBlock:     chain.ListenerStartingBlock,
		ListenerBackOffDuration:   chain.ListenerBackOffDuration,
		ListenerExternalTxEnabled: chain.ListenerExternalTxEnabled,
		Labels:                    chain.Labels,
		Headers:                   chain.Headers,
		CreatedAt:                 chain.CreatedAt,
		UpdatedAt:                 chain.UpdatedAt,
	}

	if chain.ChainID != nil {
		chainModel.ChainID = chain.ChainID.String()
	}

	if chain.PrivateTxManager != nil {
		chainModel.PrivateTxManagers = []*models.PrivateTxManager{{
			UUID:      chain.PrivateTxManager.UUID,
			ChainUUID: chain.PrivateTxManager.ChainUUID,
			URL:       chain.PrivateTxManager.URL,
			Type:      chain.PrivateTxManager.Type,
			CreatedAt: chain.PrivateTxManager.CreatedAt,
		}}
	}

	return chainModel
}
