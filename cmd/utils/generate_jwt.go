package utils

import (
	"context"
	"os"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/auth/jwt/generator"
	"github.com/ConsenSys/orchestrate/pkg/utils"
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

func run(_ *cobra.Command, _ []string) {
	// Create app
	ctx, cancel := context.WithCancel(context.Background())

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { cancel() })
	defer sig.Close()

	// Start application
	generator.Init(ctx)
	jwt, err := generator.GlobalJWTGenerator().GenerateAccessTokenWithTenantID(tenant, expiration)
	if err != nil {
		log.WithError(err).Fatalf("jwt-generator: could not generate JWT token")
	}
	log.WithFields(log.Fields{
		"jwt":      jwt,
		"tenantID": tenant,
	}).Infof("jwt-generator: token generated")
}
