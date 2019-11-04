package cucumber

import (
	"context"
	"os"
	"sync"

	"github.com/DATA-DOG/godog"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/chanregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/parser"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/steps"
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

		// Initialize parser
		parser.Init(ctx)

		// Initialize Steps
		steps.Init(ctx)

		// Initialize Channel registry
		chanregistry.Init(ctx)

		options = &godog.Options{
			ShowStepDefinitions: viper.GetBool(cucumberShowStepDefinitionsViperKey),
			Randomize:           viper.GetInt64(cucumberRandomizeViperKey),
			StopOnFailure:       viper.GetBool(cucumberStopOnFailureViperKey),
			Strict:              viper.GetBool(cucumberStrictViperKey),
			NoColors:            viper.GetBool(cucumberNoColorsViperKey),
			Tags:                viper.GetString(cucumberTagsViperKey),
			Format:              viper.GetString(cucumberFormatViperKey),
			Concurrency:         viper.GetInt(cucumberConcurrencyViperKey),
			Paths:               viper.GetStringSlice(cucumberPathsViperKey),
		}

		if viper.GetString(cucumberOutputPathViperKey) != "" {
			f, _ := os.Create(viper.GetString(cucumberOutputPathViperKey))
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
