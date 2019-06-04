package cucumber

import (
	"context"
	"os"

	"github.com/DATA-DOG/godog"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/services/cucumber/steps"
)

func Run(cancel context.CancelFunc, opt *godog.Options) {

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

	log.Info("cucumber: feature tests succeeded")
	cancel()
}
