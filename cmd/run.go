package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/config"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/utils"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	config.HTTPHostname(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, args []string) {
	// Create app
	a := app.New()

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { a.Close() })
	defer sig.Close()

	// Run App
	a.Run()

	// Wait
	<-a.Done()
}
