package parser

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt/generator"
)

var (
	parser   *Parser
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if parser != nil {
			return
		}

		// Initialize Multi-tenancy
		generator.Init(ctx)

		parser = &Parser{
			Aliases:      NewAliasRegistry(),
			JWTGenerator: generator.GlobalJWTGenerator(),
		}

		// Register aliases
		log.Infof("Registering %v aliases", len(viper.GetStringSlice(cucumberAliasesViperKey)))
		for _, alias := range viper.GetStringSlice(cucumberAliasesViperKey) {
			aka, value, err := FromAlias(alias)
			if err != nil {
				log.WithError(err).Fatalf("parser: could not parse alias")
			}
			parser.Aliases.Set("global", aka, value)
			log.WithFields(log.Fields{
				"aka":   aka,
				"value": value,
			}).Infof("parser: global alias registered")
		}
	})
}

// SetGlobalParser sets global parser
func SetGlobalParser(p *Parser) {
	parser = p
}

// GlobalParser returns global Cucumber parser
func GlobalParser() *Parser {
	return parser
}
