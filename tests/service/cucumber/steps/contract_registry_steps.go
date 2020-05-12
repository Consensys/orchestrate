package steps

import (
	"context"

	"github.com/cucumber/godog"
	gherkin "github.com/cucumber/messages-go/v10"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
)

func (sc *ScenarioContext) iRegisterTheFollowingContract(table *gherkin.PickleStepArgument_PickleTable) error {
	// Parse table
	parseContracts, err := sc.parser.ParseContracts(sc.Pickle.Id, table)
	if err != nil {
		return err
	}

	// Register parseContracts on the registry
	for _, parseContract := range parseContracts {
		_, err := sc.ContractRegistry.RegisterContract(
			authutils.WithAuthorization(context.Background(), "Bearer "+parseContract.JWTToken),
			&registry.RegisterContractRequest{
				Contract: parseContract.Contract,
			},
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func registerContractRegistrySteps(s *godog.Suite, sc *ScenarioContext) {
	s.Step(`^I register the following contract$`, sc.iRegisterTheFollowingContract)
}
