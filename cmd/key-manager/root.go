package keymanager

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "key-manager",
		Short: "Run key-manager",
	}

	rootCmd.AddCommand(newRunCommand())

	return rootCmd
}
