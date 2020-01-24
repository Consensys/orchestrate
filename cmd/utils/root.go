package utils

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "utils",
		Short: "Run utility command",
	}

	rootCmd.AddCommand(newGenerateJWTCommand())

	return rootCmd
}
