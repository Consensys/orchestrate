package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
)

var rootCmd = &cobra.Command{
	Use:              "worker",
	TraverseChildren: true,
	Version:          "v0.1.0",
}

var cmdExample = &cobra.Command{
	Use:   "example [OPTIONS]",
	Short: "An example command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Log-Level:", viper.GetString(logger.LogLevelViperKey))
		fmt.Println("Log-Format:", viper.GetString(logger.LogFormatViperKey))
	},
}

func init() {
	rootCmd.AddCommand(cmdExample)
	logger.LogLevel(cmdExample.Flags())
	logger.LogFormat(cmdExample.Flags())
}

func main() {
	_ = rootCmd.Execute()
}
