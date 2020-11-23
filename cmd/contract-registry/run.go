package contractregistry

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry"
)

var cmdErr error

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		RunE:  run,
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.PreRunBindFlags(viper.GetViper(), cmd.Flags(), "contract-registry")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			if err := errors.CombineErrors(cmdErr, cmd.Context().Err()); err != nil {
				os.Exit(1)
			}
		},
	}

	// Set flags
	contractregistry.Flags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, _ []string) error {
	return contractregistry.Run(cmd.Context())
}
