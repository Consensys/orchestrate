package controls

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/chainregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/chain-registry/utils"
)

const creditorComponent = "faucet.control.creditor"

// Controller is a controller that holds a list of account that should not be credited
type CreditorControl struct {
	chainStateReader ethclient.ChainStateReader
}

// NewController creates a new BlackList controller
func NewCreditorControl(chainStateReader ethclient.ChainStateReader) *CreditorControl {
	return &CreditorControl{
		chainStateReader: chainStateReader,
	}
}

// Control apply BlackList controller on a credit function
func (ctrl *CreditorControl) Control(ctx context.Context, req *types.Request) error {
	if len(req.Candidates) == 0 {
		return errors.FaucetWarning("no faucet candidates").ExtendComponent(creditorComponent)
	}

	for key, candidate := range req.Candidates {
		if candidate.Creditor.Hex() == req.Beneficiary.Hex() {
			delete(req.Candidates, key)
			continue
		}
		// Retrieve creditor balance
		balance, err := utils.GetAddressBalance(ctx, ctrl.chainStateReader, req.Chain.URLs, candidate.Creditor)
		if err != nil {
			return errors.FromError(err).ExtendComponent(maxBalanceComponent)
		}
		// In case balance is lower, remove candidate
		if balance.Cmp(candidate.Amount) == -1 {
			delete(req.Candidates, key)
		}
	}
	if len(req.Candidates) == 0 {
		return errors.FaucetWarning("attempt to credit the creditor").ExtendComponent(creditorComponent)
	}

	return nil
}

func (ctrl *CreditorControl) OnSelectedCandidate(_ context.Context, faucet *types.Faucet, beneficiary ethcommon.Address) error {
	return nil
}
