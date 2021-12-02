package controls

import (
	"context"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/consensys/orchestrate/pkg/types/entities"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient"
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
func (ctrl *CreditorControl) Control(ctx context.Context, req *entities.FaucetRequest) error {
	for key, candidate := range req.Candidates {
		if candidate.CreditorAccount.String() == req.Beneficiary.String() {
			delete(req.Candidates, key)
			continue
		}
		// Retrieve creditor balance
		balance, err := getAddressBalance(ctx, ctrl.chainStateReader, req.Chain.URLs, candidate.CreditorAccount)
		if err != nil {
			log.FromContext(ctx).WithError(err).Error("failed to get faucet balance")
			return errors.FromError(err).ExtendComponent(creditorComponent)
		}

		// In case balance is lower, remove candidate
		if balance.Cmp(candidate.Amount.ToInt()) == -1 {
			log.FromContext(ctx).WithField("creditor_account", candidate.CreditorAccount).
				Warn("faucet candidate discarded due to insufficient balance")

			delete(req.Candidates, key)
		}
	}

	return nil
}

func (ctrl *CreditorControl) OnSelectedCandidate(_ context.Context, _ *entities.Faucet, _ ethcommon.Address) error {
	return nil
}
