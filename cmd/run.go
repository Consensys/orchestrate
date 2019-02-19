package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/boilerplate-worker.git/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/common.git/utils"
)

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   run,
	}

	return runCmd
}

// ProcessSignals process signals
func run(cmd *cobra.Command, args []string) {
	// Create app
	a := app.New()

	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { a.Close() })
	defer sig.Close()

	// Start App
	go a.Start()

	// Wait
	<-a.Done()
	os.Exit(0)
}
