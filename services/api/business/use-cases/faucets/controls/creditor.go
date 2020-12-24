package controls

import (
	"context"
	"math/big"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"

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
	log.WithContext(ctx).Debug("creditor control check")

	for key, candidate := range req.Candidates {
		amountBigInt, _ := new(big.Int).SetString(candidate.Amount, 10)

		if candidate.CreditorAccount == req.Beneficiary {
			delete(req.Candidates, key)
			continue
		}
		// Retrieve creditor balance
		balance, err := utils.GetAddressBalance(ctx, ctrl.chainStateReader, req.Chain.URLs, candidate.CreditorAccount)
		if err != nil {
			return errors.FromError(err).ExtendComponent(creditorComponent)
		}

		// In case balance is lower, remove candidate
		if balance.Cmp(amountBigInt) == -1 {
			log.WithContext(ctx).
				WithField("balance", balance.String()).
				WithField("amount", amountBigInt.String()).
				WithField("creditor_account", candidate.CreditorAccount).
				Warning("faucet candidate discarded due to insufficient balance")

			delete(req.Candidates, key)
		}
	}

	return nil
}

func (ctrl *CreditorControl) OnSelectedCandidate(_ context.Context, _ *entities.Faucet, _ string) error {
	return nil
}
