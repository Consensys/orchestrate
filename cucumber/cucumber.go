package cucumber

import (
	"context"
	"os"

	"github.com/DATA-DOG/godog"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/cucumber/steps"
	log "github.com/sirupsen/logrus"
)

func Run(cancel context.CancelFunc, opt *godog.Options) {

	log.Info("cucumber: feature tests start")

	status := godog.RunWithOptions("tests", func(s *godog.Suite) {
		steps.FeatureContext(s)
	}, *opt)

	// godog status:
	//  0 - success
	//  1 - failed
	//  2 - command line usage error
	//  128 - or higher, os signal related error exit codes

	// If fail
	if status > 0 {
		log.Errorf("cucumber: feature tests failed with status %d", status)
		os.Exit(status)
	}

	cancel()
}
