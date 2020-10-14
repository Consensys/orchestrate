package formatters

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/identitymanager"
)

func FormatCreateIdentityRequest(req *types.CreateIdentityRequest) *entities.Identity {
	return &entities.Identity{
		Alias:      req.Alias,
		Attributes: req.Attributes,
	}
}

func FormatIdentityResponse(iden *entities.Identity) *types.IdentityResponse {
	return &types.IdentityResponse{
		Alias:               iden.Alias,
		Attributes:          iden.Attributes,
		Address:             iden.Address,
		PublicKey:           iden.PublicKey,
		CompressedPublicKey: iden.CompressedPublicKey,
		TenantID:            iden.TenantID,
		Active:              iden.Active,
		CreatedAt:           iden.CreatedAt,
		UpdatedAt:           iden.UpdatedAt,
	}
}
