package transactionscheduler

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "transaction-scheduler",
		Short: "Run transaction-scheduler",
	}

	rootCmd.AddCommand(newRunCommand())
	rootCmd.AddCommand(newMigrateCmd())

	return rootCmd
}
