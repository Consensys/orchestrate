package identitymanager

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	identitymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager"
)

var cmdErr error

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		RunE:  run,
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.PreRunBindFlags(viper.GetViper(), cmd.Flags(), "identity-manager")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			if err := errors.CombineErrors(cmdErr, cmd.Context().Err()); err != nil {
				os.Exit(1)
			}
		},
	}

	// Identity manager flags
	identitymanager.Flags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, _ []string) error {
	return identitymanager.Run(cmd.Context())
}
