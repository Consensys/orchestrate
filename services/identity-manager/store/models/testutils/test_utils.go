package testutils

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/models"
)

func FakeAccountModel() *models.Account {
	return &models.Account{
		Alias:               utils.RandomString(10),
		Active:              true,
		TenantID:            "tenantID",
		Address:             ethcommon.HexToAddress(utils.RandHexString(12)).String(),
		PublicKey:           ethcommon.HexToHash(utils.RandHexString(12)).String(),
		CompressedPublicKey: ethcommon.HexToHash(utils.RandHexString(12)).String(),
		Attributes: map[string]string{
			"attr1": "val1",
			"attr2": "val2",
		},
	}
}
