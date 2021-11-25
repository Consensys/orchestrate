package controls

import (
	"context"
	"math/big"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"

	"github.com/consensys/orchestrate/pkg/types/entities"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient"
	ethcommon "github.com/ethereum/go-ethereum/common"
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
func (ctrl *MaxBalanceControl) Control(ctx context.Context, req *entities.FaucetRequest) error {
	if len(req.Candidates) == 0 {
		return nil
	}

	// Retrieve account balance
	balance, err := getAddressBalance(ctx, ctrl.chainStateReader, req.Chain.URLs, ethcommon.HexToAddress(req.Beneficiary))
	if err != nil {
		log.FromContext(ctx).WithError(err).Error("failed to get faucet balance")
		return errors.FromError(err).ExtendComponent(maxBalanceComponent)
	}

	// Ensure MaxBalance is respected
	for key, candidate := range req.Candidates {
		amountBigInt, _ := new(big.Int).SetString(candidate.Amount, 10)
		maxBalanceBigInt, _ := new(big.Int).SetString(candidate.MaxBalance, 10)

		if new(big.Int).Add(amountBigInt, balance).Cmp(maxBalanceBigInt) > 0 {
			delete(req.Candidates, key)
		}
	}

	return nil
}

func (ctrl *MaxBalanceControl) OnSelectedCandidate(_ context.Context, _ *entities.Faucet, _ string) error {
	return nil
}
