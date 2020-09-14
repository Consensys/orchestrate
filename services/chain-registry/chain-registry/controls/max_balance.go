package controls

import (
	"context"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/chainregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/utils"
)

const maxBalanceComponent = "faucet.control.max-balance"

// Controller is a controller that ensures an address can not be credit above a given limit
type MaxBalanceControl struct {
	chainStateReader ethclient.ChainStateReader
}

// NewController creates a new max balance controller
func NewMaxBalanceControl(chainStateReader ethclient.ChainStateReader) *MaxBalanceControl {
	return &MaxBalanceControl{
		chainStateReader: chainStateReader,
	}
}

// Control apply MaxBalance controller on a credit function
func (ctrl *MaxBalanceControl) Control(ctx context.Context, req *types.Request) error {
	if len(req.Candidates) == 0 {
		return errors.FaucetWarning("no faucet candidates").ExtendComponent(maxBalanceComponent)
	}

	// Retrieve account balance
	balance, err := utils.GetAddressBalance(ctx, ctrl.chainStateReader, req.Chain.URLs, req.Beneficiary)
	if err != nil {
		return errors.FromError(err).ExtendComponent(maxBalanceComponent)
	}

	// Ensure MaxBalance is respected
	for key, candidate := range req.Candidates {
		if new(big.Int).Add(candidate.Amount, balance).Cmp(candidate.MaxBalance) > 0 {
			delete(req.Candidates, key)
		}
	}
	if len(req.Candidates) == 0 {
		// Do not credit if final balance would be superior to max authorized
		return errors.FaucetWarning("account balance too high").ExtendComponent(maxBalanceComponent)
	}

	return nil
}

func (ctrl *MaxBalanceControl) OnSelectedCandidate(_ context.Context, _ *types.Faucet, _ ethcommon.Address) error {
	return nil
}
