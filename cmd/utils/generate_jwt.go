package utils

import (
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/jwt/generator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	tenant     string
	expiration time.Duration
)

func newGenerateJWTCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "generate-jwt",
		Short: "Generate JWT Access Token",
		Run:   run,
	}

	generator.PrivateKey(runCmd.Flags())
	runCmd.Flags().StringVar(&tenant, "tenant", "_", "Tenant to create a token for")
	runCmd.Flags().DurationVar(&expiration, "expiration", time.Hour, "Token expiration time")

	return runCmd
}

func run(cmd *cobra.Command, _ []string) {
	// Start application
	generator.Init(cmd.Context())

	gJwt := generator.GlobalJWTGenerator()
	if gJwt == nil {
		log.Fatal("jwt-generator: could not initialize it")
	}

	jwt, err := gJwt.GenerateAccessTokenWithTenantID(tenant, []string{"*:*"}, expiration)
	if err != nil {
		log.WithError(err).Fatal("jwt-generator: could not generate JWT token")
	}
	log.WithFields(log.Fields{
		"jwt":      jwt,
		"tenantID": tenant,
	}).Info("jwt-generator: token generated")
}
