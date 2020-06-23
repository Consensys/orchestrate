package steps

import (
	"context"

	"github.com/cucumber/godog"
	gherkin "github.com/cucumber/messages-go/v10"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/utils"
)

func (sc *ScenarioContext) iRegisterTheFollowingContract(table *gherkin.PickleStepArgument_PickleTable) error {
	// Parse table
	parseContracts, err := utils.ParseContracts(table)
	if err != nil {
		return err
	}

	// Register parseContracts on the registry
	for _, parseContract := range parseContracts {
		_, err := sc.ContractRegistry.RegisterContract(
			authutils.WithAuthorization(context.Background(), parseContract.JWTToken),
			&registry.RegisterContractRequest{
				Contract: parseContract.Contract,
			},
		)

		if err != nil {
			return err
		}
		sc.TearDownFunc = append(sc.TearDownFunc, func() {
			_, _ = sc.ContractRegistry.DeregisterContract(
				authutils.WithAuthorization(context.Background(), parseContract.JWTToken),
				&registry.DeregisterContractRequest{
					ContractId: parseContract.Contract.Id,
				},
			)
		})
	}

	return nil
}

func registerContractRegistrySteps(s *godog.ScenarioContext, sc *ScenarioContext) {
	s.Step(`^I register the following contracts$`, sc.preProcessTableStep(sc.iRegisterTheFollowingContract))
}
