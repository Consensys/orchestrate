package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/tests/handlers/logger"
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
		fmt.Println("Log-Level:", viper.GetString("log.level"))
		fmt.Println("Log-Format:", viper.GetString("log.format"))
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
