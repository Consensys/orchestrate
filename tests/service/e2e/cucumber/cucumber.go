package cucumber

import (
	"github.com/cucumber/godog"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/cucumber/steps"
)

func Run(opt *godog.Options) error {
	status := godog.TestSuite{
		Name:                "tests",
		ScenarioInitializer: steps.InitializeScenario,
		Options:             opt,
	}.Run()

	// godog status:
	//  0 - success
	//  1 - failed
	//  2 - command line usage error
	//  128 - or higher, os signal related error exit codes

	// If fail
	if status > 0 {
		return errors.InternalError("cucumber: feature tests failed with status %d", status)
	}

	log.Info("cucumber: feature tests succeeded")
	return nil
}
