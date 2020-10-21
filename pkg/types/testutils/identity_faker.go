package testutils

import (
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func FakeAccount() *entities.Account {
	return &entities.Account{
		Alias:               uuid.Must(uuid.NewV4()).String(),
		TenantID:            utils.RandomString(6),
		Attributes:          make(map[string]string),
		Address:             ethcommon.HexToAddress(utils.RandHexString(12)).String(),
		PublicKey:           ethcommon.HexToHash(utils.RandHexString(12)).String(),
		CompressedPublicKey: ethcommon.HexToHash(utils.RandHexString(12)).String(),
		Active:              true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}
