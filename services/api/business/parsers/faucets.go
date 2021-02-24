package parsers

import (
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/api/store/models"
)

func NewFaucetFromModel(faucet *models.Faucet) *entities.Faucet {
	return &entities.Faucet{
		UUID:            faucet.UUID,
		Name:            faucet.Name,
		TenantID:        faucet.TenantID,
		ChainRule:       faucet.ChainRule,
		CreditorAccount: faucet.CreditorAccount,
		MaxBalance:      faucet.MaxBalance,
		Amount:          faucet.Amount,
		Cooldown:        faucet.Cooldown,
		CreatedAt:       faucet.CreatedAt,
		UpdatedAt:       faucet.UpdatedAt,
	}
}

func NewFaucetModelFromEntity(faucet *entities.Faucet) *models.Faucet {
	return &models.Faucet{
		UUID:            faucet.UUID,
		Name:            faucet.Name,
		TenantID:        faucet.TenantID,
		ChainRule:       faucet.ChainRule,
		CreditorAccount: faucet.CreditorAccount,
		MaxBalance:      faucet.MaxBalance,
		Amount:          faucet.Amount,
		Cooldown:        faucet.Cooldown,
		CreatedAt:       faucet.CreatedAt,
	}
}
