package parsers

import (
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/chainregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
)

func NewFaucetModelFromEntity(faucet *types.Faucet) *models.Faucet {
	return &models.Faucet{
		UUID:            faucet.UUID,
		CreditorAccount: faucet.Creditor.Hex(),
		MaxBalance:      faucet.MaxBalance.String(),
		Cooldown:        faucet.Cooldown.String(),
		Amount:          faucet.Amount.String(),
	}
}

func NewFaucetEntitiesFromModels(faucets []*models.Faucet) map[string]types.Faucet {
	eligibleFaucets := make(map[string]types.Faucet)

	for _, f := range faucets {
		eligibleFaucets[f.UUID] = types.Faucet{
			UUID:       f.UUID,
			Creditor:   ethcommon.HexToAddress(f.CreditorAccount),
			MaxBalance: func() (b *big.Int) { b, _ = new(big.Int).SetString(f.MaxBalance, 10); return }(),
			Amount:     func() (a *big.Int) { a, _ = new(big.Int).SetString(f.Amount, 10); return }(),
			Cooldown:   func() (d time.Duration) { d, _ = time.ParseDuration(f.Cooldown); return }(),
		}
	}

	return eligibleFaucets
}
