package steps

import (
	"context"

	"encoding/json"

	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/tests/service/e2e/utils"
	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/cucumber/godog"
	gherkin "github.com/cucumber/messages-go/v10"
	"github.com/traefik/traefik/v2/pkg/log"
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

		headers := utils.GetHeaders(parseContract.APIKey, parseContract.Tenant, parseContract.JWTToken)
		_, err = sc.client.RegisterContract(
			context.WithValue(ctx, client.RequestHeaderKey, headers),
			&api.RegisterContractRequest{
				Name:             parseContract.Contract.Name,
				Tag:              parseContract.Contract.Tag,
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
				WithField("name", parseContract.Contract.Name).
				WithField("tag", parseContract.Contract.Tag).
				Warn("DeregisterContract is not implemented")
		})
	}

	return nil
}

func initContractRegistrySteps(s *godog.ScenarioContext, sc *ScenarioContext) {
	s.Step(`^I register the following contracts$`, sc.preProcessTableStep(sc.iRegisterTheFollowingContract))
}
