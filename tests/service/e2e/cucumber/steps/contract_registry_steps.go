package steps

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/cucumber/godog"
	gherkin "github.com/cucumber/messages-go/v10"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/utils"
)

func (sc *ScenarioContext) iRegisterTheFollowingContract(table *gherkin.PickleStepArgument_PickleTable) error {
	ctx := context.Background()

	// Parse table
	parseContracts, err := utils.ParseContracts(table)
	if err != nil {
		return err
	}

	// Register parseContracts on the registry
	for _, parseContract := range parseContracts {
		var abi interface{}
		err := json.Unmarshal([]byte(parseContract.Contract.ABI), &abi)
		if err != nil {
			return err
		}

		_, err = sc.client.RegisterContract(
			authutils.WithAuthorization(ctx, parseContract.JWTToken),
			&api.RegisterContractRequest{
				Name:             parseContract.Contract.ID.Name,
				Tag:              parseContract.Contract.ID.Tag,
				ABI:              abi,
				Bytecode:         parseContract.Contract.Bytecode,
				DeployedBytecode: parseContract.Contract.DeployedBytecode,
			},
		)

		if err != nil {
			return err
		}

		sc.TearDownFunc = append(sc.TearDownFunc, func() {
			log.FromContext(ctx).
				WithField("name", parseContract.Contract.ID.Name).
				WithField("tag", parseContract.Contract.ID.Tag).
				Warn("DeregisterContract is not implemented")
		})
	}

	return nil
}

func registerContractRegistrySteps(s *godog.ScenarioContext, sc *ScenarioContext) {
	s.Step(`^I register the following contracts$`, sc.preProcessTableStep(sc.iRegisterTheFollowingContract))
}
