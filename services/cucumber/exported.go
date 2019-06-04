package cucumber

import (
	"context"
	"os"
	"sync"

	"github.com/DATA-DOG/godog"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/services/chanregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/services/cucumber/steps"
)

var (
	options  *godog.Options
	initOnce = &sync.Once{}
)

// Init initialize Cucumber service
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if options != nil {
			return
		}

		// Initialize Steps
		steps.Init(ctx)

		// Initialize Channel registry
		chanregistry.Init(ctx)

		options = &godog.Options{
			ShowStepDefinitions: viper.GetBool("cucumber.showstepdefinitions"),
			Randomize:           viper.GetInt64("cucumber.randomize"),
			StopOnFailure:       viper.GetBool("cucumber.stoponfailure"),
			Strict:              viper.GetBool("cucumber.strict"),
			NoColors:            viper.GetBool("cucumber.nocolors"),
			Tags:                viper.GetString("cucumber.tags"),
			Format:              viper.GetString("cucumber.format"),
			Concurrency:         viper.GetInt("cucumber.concurrency"),
			Paths:               viper.GetStringSlice("cucumber.paths"),
		}

		if viper.GetString("cucumber.outputpath") != "" {
			f, _ := os.Create(viper.GetString("cucumber.outputpath"))
			options.Output = f
		}

		log.Infof("cucumber: service ready")
	})
}

// SetGlobalOptions sets global Cucumber Handler
func SetGlobalOptions(o *godog.Options) {
	options = o
}

// GlobalHandler returns global Cucumber handler
func GlobalOptions() *godog.Options {
	return options
}
