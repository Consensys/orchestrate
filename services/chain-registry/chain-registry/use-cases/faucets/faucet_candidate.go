package faucets

import (
	"context"
	"reflect"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/chain-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/controls"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/use-cases/chains"
)

const faucetCandidateComponent = "use-cases.faucet-candidate"

type FaucetCandidate interface {
	Execute(context.Context, ethcommon.Address, string, []string) (*types.Faucet, error)
}

type FaucetControl interface {
	Control(context.Context, *types.Request) error
	OnSelectedCandidate(context.Context, *types.Faucet, ethcommon.Address) error
}

// RegisterContract is a use case to register a new contract
type faucetCandidate struct {
	getChainUC       chains.GetChain
	chainStateReader ethclient.ChainStateReader
	getFaucetsUC     GetFaucets
	controls         []FaucetControl
}

// NewGetCatalog creates a new GetCatalog
func NewFaucetCandidateUseCase(getChainUC chains.GetChain, getFaucets GetFaucets, chainStateReader ethclient.ChainStateReader) FaucetCandidate {
	cooldownCtrl := controls.NewCooldownControl()
	maxBalanceCtrl := controls.NewMaxBalanceControl(chainStateReader)
	creditorCtrl := controls.NewCreditorControl(chainStateReader)

	return &faucetCandidate{
		getChainUC:       getChainUC,
		chainStateReader: chainStateReader,
		getFaucetsUC:     getFaucets,
		controls:         []FaucetControl{creditorCtrl, cooldownCtrl, maxBalanceCtrl},
	}
}

func (uc *faucetCandidate) Execute(ctx context.Context, account ethcommon.Address, chainUUID string, tenants []string) (*types.Faucet, error) {
	faucets, err := uc.getFaucetsUC.Execute(ctx, []string{}, map[string]string{"chain_rule": chainUUID})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(faucetCandidateComponent)
	}

	chain, err := uc.getChainUC.Execute(ctx, chainUUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(faucetCandidateComponent)
	}

	req := &types.Request{
		Beneficiary: account,
		Candidates:  parsers.NewFaucetEntitiesFromModels(faucets),
		Chain:       chain,
	}

	for _, ctrl := range uc.controls {
		err = ctrl.Control(ctx, req)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(faucetCandidateComponent)
		}
	}

	if len(req.Candidates) < 1 {
		return nil, nil
	}

	// Select a first faucet candidate for comparison
	selectedFaucet := req.Candidates[electFaucet(req.Candidates)]

	for _, ctrl := range uc.controls {
		err := ctrl.OnSelectedCandidate(ctx, &selectedFaucet, req.Beneficiary)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(faucetCandidateComponent)
		}
	}

	return &selectedFaucet, nil
}

// electFaucet is currently selecting the remaining faucet candidates with the highest amount
func electFaucet(faucetsCandidates map[string]types.Faucet) string {
	// Select a first faucet candidate for comparison
	electedFaucet := reflect.ValueOf(faucetsCandidates).MapKeys()[0].String()
	for key, candidate := range faucetsCandidates {
		if candidate.Amount.Cmp(faucetsCandidates[electedFaucet].Amount) > 0 {
			electedFaucet = key
		}
	}
	return electedFaucet
}
