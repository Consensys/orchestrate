package cucumber

import (
	"context"
	"os"
	"strings"
	"sync"

	broker "github.com/consensys/orchestrate/pkg/broker/sarama"
	"github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	ethclient "github.com/consensys/orchestrate/pkg/toolkit/ethclient/rpc"
	"github.com/consensys/orchestrate/tests/service/e2e/cucumber/alias"
	"github.com/cucumber/godog"
	"github.com/spf13/viper"
)

var (
	options  *godog.Options
	initOnce = &sync.Once{}
)

// Init initialize Cucumber service
func Init(ctx context.Context, rawTestData string) {
	initOnce.Do(func() {
		if options != nil {
			return
		}

		logger := log.FromContext(ctx)

		// Initialize Steps
		broker.InitSyncProducer(ctx)
		alias.Init(rawTestData)
		client.Init()
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

		outputPath := viper.GetString(OutputPathViperKey)
		if outputPath != "" {
			f, err := os.Create(viper.GetString(OutputPathViperKey))
			if err != nil {
				logger.WithError(err).Fatalf("could not write output in %s", outputPath)
			}
			options.Output = f
		}

		logger.WithField("tags", options.Tags).
			WithField("concurrency", options.Concurrency).
			WithField("paths", options.Paths).
			WithField("output", outputPath).
			Info("service ready")
	})
}

func listTagCucumber() string {
	var tags []string
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

// GlobalOptions returns global Cucumber handler
func GlobalOptions() *godog.Options {
	return options
}
