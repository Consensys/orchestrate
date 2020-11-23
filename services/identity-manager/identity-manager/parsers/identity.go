package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/store/models"
)

func NewAccountModelFromEntities(iden *entities.Account) *models.Account {
	return &models.Account{
		Alias:               iden.Alias,
		Address:             iden.Address,
		PublicKey:           iden.PublicKey,
		CompressedPublicKey: iden.CompressedPublicKey,
		TenantID:            iden.TenantID,
		Attributes:          iden.Attributes,
		CreatedAt:           iden.CreatedAt,
		UpdatedAt:           iden.UpdatedAt,
	}
}

func NewAccountEntityFromModels(iden *models.Account) *entities.Account {
	return &entities.Account{
		Alias:               iden.Alias,
		Address:             iden.Address,
		PublicKey:           iden.PublicKey,
		CompressedPublicKey: iden.CompressedPublicKey,
		TenantID:            iden.TenantID,
		Attributes:          iden.Attributes,
		CreatedAt:           iden.CreatedAt,
		UpdatedAt:           iden.UpdatedAt,
	}
}
