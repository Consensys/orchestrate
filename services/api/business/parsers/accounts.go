package parsers

import (
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/store/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func NewAccountModelFromEntities(account *entities.Account) *models.Account {
	return &models.Account{
		Alias:               account.Alias,
		Address:             account.Address.Hex(),
		PublicKey:           account.PublicKey,
		CompressedPublicKey: account.CompressedPublicKey,
		TenantID:            account.TenantID,
		OwnerID:             account.OwnerID,
		Attributes:          account.Attributes,
		CreatedAt:           account.CreatedAt,
		UpdatedAt:           account.UpdatedAt,
	}
}

func NewAccountEntityFromModels(account *models.Account) *entities.Account {
	return &entities.Account{
		Alias:               account.Alias,
		Address:             ethcommon.HexToAddress(account.Address),
		PublicKey:           account.PublicKey,
		CompressedPublicKey: account.CompressedPublicKey,
		TenantID:            account.TenantID,
		OwnerID:             account.OwnerID,
		Attributes:          account.Attributes,
		CreatedAt:           account.CreatedAt,
		UpdatedAt:           account.UpdatedAt,
	}
}
