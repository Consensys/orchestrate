package steps

import (
	"context"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"

	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
)

func (sc *ScenarioContext) iRegisterTheFollowingContract(table *gherkin.DataTable) error {
	// Parse table
	parseContracts, err := sc.parser.ParseContracts(sc.ID, table)
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
