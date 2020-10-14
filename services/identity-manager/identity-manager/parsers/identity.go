package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/models"
)

func NewIdentityModelFromEntities(iden *entities.Identity) *models.Identity {
	return &models.Identity{
		Alias:               iden.Alias,
		Address:             iden.Address,
		PublicKey:           iden.PublicKey,
		CompressedPublicKey: iden.CompressedPublicKey,
		TenantID:            iden.TenantID,
		Active:              iden.Active,
		Attributes:          iden.Attributes,
		CreatedAt:           iden.CreatedAt,
		UpdatedAt:           iden.UpdatedAt,
	}
}

func NewIdentityEntityFromModels(iden *models.Identity) *entities.Identity {
	return &entities.Identity{
		Alias:               iden.Alias,
		Address:             iden.Address,
		PublicKey:           iden.PublicKey,
		CompressedPublicKey: iden.CompressedPublicKey,
		TenantID:            iden.TenantID,
		Active:              iden.Active,
		Attributes:          iden.Attributes,
		CreatedAt:           iden.CreatedAt,
		UpdatedAt:           iden.UpdatedAt,
	}
}
