package types

import (
	"math/big"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// Envelope holds information for a Faucet Credit Envelope
type Request struct {
	ParentTxID        string
	ChildTxID         string
	ChainID           *big.Int
	ChainURL          string
	ChainName         string
	ChainUUID         string
	Beneficiary       ethcommon.Address
	FaucetsCandidates map[string]Faucet
	ElectedFaucet     string
}

type Faucet struct {
	Amount     *big.Int
	MaxBalance *big.Int
	Cooldown   time.Duration
	Creditor   ethcommon.Address
}

func NewFaucetsCandidates(storeFaucet []*models.Faucet) map[string]Faucet {
	eligibleFaucets := make(map[string]Faucet)

	for _, f := range storeFaucet {
		eligibleFaucets[f.UUID] = Faucet{
			Creditor:   ethcommon.HexToAddress(f.CreditorAccount),
			MaxBalance: func() (b *big.Int) { b, _ = new(big.Int).SetString(f.MaxBalance, 10); return }(),
			Amount:     func() (a *big.Int) { a, _ = new(big.Int).SetString(f.Amount, 10); return }(),
			Cooldown:   func() (d time.Duration) { d, _ = time.ParseDuration(f.Cooldown); return }(),
		}
	}

	return eligibleFaucets
}
