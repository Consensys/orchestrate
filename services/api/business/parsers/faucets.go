package parsers

import (
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/store/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func NewFaucetFromModel(faucet *models.Faucet) *entities.Faucet {
	return &entities.Faucet{
		UUID:            faucet.UUID,
		Name:            faucet.Name,
		TenantID:        faucet.TenantID,
		ChainRule:       faucet.ChainRule,
		CreditorAccount: ethcommon.HexToAddress(faucet.CreditorAccount),
		MaxBalance:      faucet.MaxBalance,
		Amount:          faucet.Amount,
		Cooldown:        faucet.Cooldown,
		CreatedAt:       faucet.CreatedAt,
		UpdatedAt:       faucet.UpdatedAt,
	}
}

func NewFaucetModelFromEntity(faucet *entities.Faucet) *models.Faucet {
	f := &models.Faucet{
		UUID:            faucet.UUID,
		Name:            faucet.Name,
		TenantID:        faucet.TenantID,
		ChainRule:       faucet.ChainRule,
		CreditorAccount: faucet.CreditorAccount.Hex(),
		MaxBalance:      faucet.MaxBalance,
		Amount:          faucet.Amount,
		Cooldown:        faucet.Cooldown,
		CreatedAt:       faucet.CreatedAt,
	}

	return f
}
