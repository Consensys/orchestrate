package cucumber

import (
	"context"
	"os"
	"sync"

	"github.com/cucumber/godog"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

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

		tags := listTagCucumber()

		options = &godog.Options{
			ShowStepDefinitions: viper.GetBool(ShowStepDefinitionsViperKey),
			Randomize:           viper.GetInt64(RandomizeViperKey),
			StopOnFailure:       viper.GetBool(StopOnFailureViperKey),
			Strict:              viper.GetBool(StrictViperKey),
			NoColors:            viper.GetBool(NoColorsViperKey),
			Tags:                viper.GetString(TagsViperKey),
			Format:              viper.GetString(FormatViperKey),
			Concurrency:         viper.GetInt(ConcurrencyViperKey),
			Paths:               viper.GetStringSlice(PathsViperKey),
		}

		if viper.GetString(OutputPathViperKey) != "" {
			f, _ := os.Create(viper.GetString(OutputPathViperKey))
			options.Output = f
		}

		log.Infof("cucumber: service ready")
	})
}

func listTagCucumber() (tags string) {
	externalTag := viper.GetString(TagsViperKey)
	if externalTag != "" && !viper.GetBool(multitenancy.EnabledViperKey) {
		tags = externalTag + " && ~@multi-tenancy"
	} else if !viper.GetBool(multitenancy.EnabledViperKey) {
		tags = " ~@multi-tenancy"
	}

	return tags
}

// SetGlobalOptions sets global Cucumber Handler
func SetGlobalOptions(o *godog.Options) {
	options = o
}

// GlobalHandler returns global Cucumber handler
func GlobalOptions() *godog.Options {
	return options
}
