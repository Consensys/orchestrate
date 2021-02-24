package parsers

import (
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/api/store/models"
)

func NewAccountModelFromEntities(account *entities.Account) *models.Account {
	return &models.Account{
		Alias:               account.Alias,
		Address:             account.Address,
		PublicKey:           account.PublicKey,
		CompressedPublicKey: account.CompressedPublicKey,
		TenantID:            account.TenantID,
		Attributes:          account.Attributes,
		CreatedAt:           account.CreatedAt,
		UpdatedAt:           account.UpdatedAt,
	}
}

func NewAccountEntityFromModels(account *models.Account) *entities.Account {
	return &entities.Account{
		Alias:               account.Alias,
		Address:             account.Address,
		PublicKey:           account.PublicKey,
		CompressedPublicKey: account.CompressedPublicKey,
		TenantID:            account.TenantID,
		Attributes:          account.Attributes,
		CreatedAt:           account.CreatedAt,
		UpdatedAt:           account.UpdatedAt,
	}
}
