package cucumber

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/cucumber/godog"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt/generator"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"
	identitymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/client"
	noncememory "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce/memory"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/e2e/cucumber/alias"
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
		broker.InitSyncProducer(ctx)
		generator.Init(ctx)
		chainregistry.Init(ctx)
		alias.Init(ctx)
		contractregistry.Init(ctx)
		txscheduler.Init()
		identitymanager.Init()
		noncememory.Init(ctx)
		ethclient.Init(ctx)

		tags := listTagCucumber()

		options = &godog.Options{
			ShowStepDefinitions: viper.GetBool(ShowStepDefinitionsViperKey),
			Randomize:           viper.GetInt64(RandomizeViperKey),
			StopOnFailure:       viper.GetBool(StopOnFailureViperKey),
			Strict:              viper.GetBool(StrictViperKey),
			NoColors:            viper.GetBool(NoColorsViperKey),
			Tags:                tags,
			Format:              viper.GetString(FormatViperKey),
			Concurrency:         viper.GetInt(ConcurrencyViperKey),
			Paths:               viper.GetStringSlice(PathsViperKey),
		}

		if outputPath := viper.GetString(OutputPathViperKey); outputPath != "" {
			f, err := os.Create(viper.GetString(OutputPathViperKey))
			if err != nil {
				log.WithError(err).Fatalf("cucumber: could not write output in %s", outputPath)
			}
			options.Output = f
		}

		log.WithField("tags", options.Tags).
			WithField("concurrency", options.Concurrency).
			WithField("paths", options.Paths).
			Infof("cucumber: service ready")
	})
}

func listTagCucumber() string {
	tags := []string{}
	if viper.GetString(TagsViperKey) != "" {
		tags = append(tags, strings.Split(viper.GetString(TagsViperKey), " ")...)
	}

	if !viper.GetBool(multitenancy.EnabledViperKey) {
		tags = append(tags, "~@multi-tenancy")
	}

	return strings.Join(tags, " && ")
}

// SetGlobalOptions sets global Cucumber Handler
func SetGlobalOptions(o *godog.Options) {
	options = o
}

// GlobalHandler returns global Cucumber handler
func GlobalOptions() *godog.Options {
	return options
}
